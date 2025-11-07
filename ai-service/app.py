import os
import pickle
import numpy as np
from flask import Flask, request, jsonify
from PIL import Image
import tensorflow as tf
# Use tf_keras for compatibility with legacy Keras 2.x models
import tf_keras as keras

app = Flask(__name__)

class ChickenDiseaseDetector:
    def __init__(self, model_path, label_encoder_path):
        """Initialize the disease detector with model and label encoder"""
        self.model = None
        self.label_encoder = None
        self.input_shape = (224, 224)  # Default input shape, can be updated based on model
        
        self.load_model(model_path)
        self.load_label_encoder(label_encoder_path)
    
    def load_model(self, model_path):
        """Load the trained Keras model"""
        try:
            self.model = keras.models.load_model(model_path)
            print(f"Model loaded successfully from {model_path}")
            print(f"Model input shape: {self.model.input_shape}")
            # Update input shape based on model
            if len(self.model.input_shape) >= 3:
                self.input_shape = self.model.input_shape[1:3]
        except Exception as e:
            print(f"Error loading model: {e}")
            raise
    
    def load_label_encoder(self, label_encoder_path):
        """Load the label encoder"""
        try:
            with open(label_encoder_path, 'rb') as f:
                self.label_encoder = pickle.load(f)
            print(f"Label encoder loaded successfully from {label_encoder_path}")
            print(f"Classes: {self.label_encoder.classes_}")
        except Exception as e:
            print(f"Error loading label encoder: {e}")
            raise
    
    def preprocess_image(self, image):
        """Preprocess image for model prediction"""
        try:
            # Convert to RGB if needed
            if image.mode != 'RGB':
                image = image.convert('RGB')
            
            # Resize to model input size
            image = image.resize(self.input_shape)
            
            # Convert to numpy array
            img_array = np.array(image)
            
            # Normalize pixel values (adjust based on training preprocessing)
            img_array = img_array.astype('float32') / 255.0
            
            # Add batch dimension
            img_array = np.expand_dims(img_array, axis=0)
            
            return img_array
        except Exception as e:
            print(f"Error preprocessing image: {e}")
            raise
    
    def predict_disease(self, image):
        """Predict disease from preprocessed image"""
        try:
            # Preprocess the image
            processed_image = self.preprocess_image(image)
            
            # Make prediction with error handling
            print("Starting model prediction...")
            try:
                predictions = self.model.predict(processed_image, verbose=0)
            except Exception as pred_error:
                print(f"Prediction failed, retrying... Error: {pred_error}")
                # Retry once
                predictions = self.model.predict(processed_image, verbose=0)
            
            # Get predicted class
            predicted_class_idx = np.argmax(predictions[0])
            confidence = float(predictions[0][predicted_class_idx])
            
            # Convert to label
            predicted_label = self.label_encoder.inverse_transform([predicted_class_idx])[0]
            
            # Get all probabilities
            all_probabilities = {}
            for i, prob in enumerate(predictions[0]):
                label = self.label_encoder.inverse_transform([i])[0]
                all_probabilities[label] = float(prob)
            
            return {
                'predicted_disease': predicted_label,
                'confidence': confidence,
                'all_probabilities': all_probabilities,
                'is_healthy': predicted_label.lower() in ['healthy', 'normal', 'no_disease'],
                'recommendation': self.get_recommendation(predicted_label, confidence)
            }
        except Exception as e:
            print(f"Error during prediction: {e}")
            import traceback
            traceback.print_exc()
            raise
    
    def get_recommendation(self, disease, confidence):
        """Get recommendation based on predicted disease"""
        recommendations = {
            'healthy': "The chicken appears healthy. Continue regular monitoring.",
            'coccidiosis': "Possible coccidiosis detected. Consider consulting a veterinarian and check water quality.",
            'salmonella': "Potential salmonella infection. Isolate the bird and consult a veterinarian immediately.",
            'e_coli': "Possible E.coli infection. Improve hygiene and consult a veterinarian.",
            'newcastle': "Potential Newcastle disease. This is serious - contact veterinarian immediately and isolate birds.",
        }
        
        base_recommendation = recommendations.get(disease.lower(), "Unknown condition detected. Consult a veterinarian for proper diagnosis.")
        
        if confidence < 0.7:
            base_recommendation += " Note: Low confidence prediction - consider retaking photo with better lighting and closer view."
        
        return base_recommendation

# Initialize the detector (update paths to your model files)
import os
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
MODEL_PATH = os.path.join(BASE_DIR, "model", "chicken_disease_model_efficientnetb0_final.h5")
LABEL_ENCODER_PATH = os.path.join(BASE_DIR, "model", "label_encoder.pkl")

detector = None

def initialize_detector():
    """Initialize detector with error handling"""
    global detector
    try:
        print(f"Checking paths:")
        print(f"  Model path: {MODEL_PATH}")
        print(f"  Encoder path: {LABEL_ENCODER_PATH}")
        print(f"  Model exists: {os.path.exists(MODEL_PATH)}")
        print(f"  Encoder exists: {os.path.exists(LABEL_ENCODER_PATH)}")
        
        if os.path.exists(MODEL_PATH) and os.path.exists(LABEL_ENCODER_PATH):
            print("Creating ChickenDiseaseDetector...")
            detector = ChickenDiseaseDetector(MODEL_PATH, LABEL_ENCODER_PATH)
            print("Disease detector initialized successfully!")
            print(f"Detector is now: {detector}")
        else:
            print(f"Model files not found:")
            print(f"  Model: {MODEL_PATH} - {'Found' if os.path.exists(MODEL_PATH) else 'Not found'}")
            print(f"  Label Encoder: {LABEL_ENCODER_PATH} - {'Found' if os.path.exists(LABEL_ENCODER_PATH) else 'Not found'}")
    except Exception as e:
        print(f"Failed to initialize detector: {e}")
        import traceback
        traceback.print_exc()

@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint"""
    global detector
    print(f"Health check - detector is: {detector}")
    print(f"Health check - detector is None: {detector is None}")
    return jsonify({
        'status': 'healthy' if detector else 'model_not_loaded',
        'model_loaded': detector is not None,
        'detector_type': str(type(detector)) if detector else 'None'
    })

@app.route('/predict', methods=['POST'])
def predict_disease():
    """Main prediction endpoint"""
    if not detector:
        print("ERROR: Detector not initialized")
        return jsonify({
            'success': False,
            'error': 'Model not loaded. Please check server logs.',
            'details': 'AI model initialization failed'
        }), 500
    
    try:
        print(f"Received prediction request")
        print(f"Request files: {request.files}")
        print(f"Request form: {request.form}")
        print(f"Request content type: {request.content_type}")
        
        # Check if image is provided
        if 'image' not in request.files:
            print("ERROR: No 'image' field in request.files")
            return jsonify({
                'success': False,
                'error': 'No image provided',
                'hint': 'Expected multipart/form-data with image field'
            }), 400
        
        file = request.files['image']
        if file.filename == '':
            print("ERROR: Empty filename")
            return jsonify({
                'success': False,
                'error': 'No image selected'
            }), 400
        
        print(f"Processing image: {file.filename}, content_type: {file.content_type}")
        
        # Read and process image
        try:
            image = Image.open(file.stream)
            print(f"Image opened successfully: size={image.size}, mode={image.mode}")
        except Exception as img_error:
            print(f"ERROR: Failed to open image: {img_error}")
            return jsonify({
                'success': False,
                'error': 'Invalid image file',
                'details': str(img_error)
            }), 400
        
        # Make prediction
        print("Starting prediction...")
        try:
            result = detector.predict_disease(image)
            print(f"Prediction complete: {result['predicted_disease']} ({result['confidence']:.2%})")
            
            return jsonify({
                'success': True,
                'prediction': result,
                'timestamp': str(np.datetime64('now'))
            })
        except Exception as pred_error:
            print(f"ERROR: Prediction failed: {pred_error}")
            import traceback
            traceback.print_exc()
            return jsonify({
                'success': False,
                'error': 'Prediction failed',
                'details': str(pred_error)
            }), 500
    
    except Exception as e:
        print(f"ERROR: Unexpected error: {e}")
        import traceback
        traceback.print_exc()
        return jsonify({
            'success': False,
            'error': 'Internal server error',
            'details': str(e),
            'type': str(type(e).__name__)
        }), 500

if __name__ == '__main__':
    # Initialize the detector
    print("="*60)
    print("Starting Tokkatot AI Service")
    print("="*60)
    initialize_detector()
    
    if detector:
        print("✅ AI Service ready to accept requests")
    else:
        print("❌ AI Service started but detector failed to initialize")
    
    print("="*60)
    
    # Start the Flask server on localhost (middleware will proxy requests)
    print("Starting Flask server on 127.0.0.1:5000")
    app.run(host='127.0.0.1', port=5000, debug=False)

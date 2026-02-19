"""
Inference module for ensemble chicken disease detection.
Provides easy-to-use interface for prediction with safety-first logic.
"""

import torch
from PIL import Image
from pathlib import Path
import numpy as np
from typing import Union, Dict, Tuple

from models import create_ensemble
from data_utils import get_transforms, CLASS_NAMES, IDX_TO_CLASS


class ChickenDiseaseDetector:
    """
    High-level interface for chicken disease detection with safety-first ensemble.
    """
    
    def __init__(
        self,
        model_path: str,
        device: str = 'auto',
        healthy_threshold: float = 0.80,
        uncertainty_threshold: float = 0.50
    ):
        """
        Initialize the detector.
        
        Args:
            model_path: Path to saved ensemble model
            device: 'cuda', 'cpu', or 'auto'
            healthy_threshold: Confidence threshold for healthy classification (0.80 = 80%)
            uncertainty_threshold: Minimum confidence threshold for any prediction (0.50 = 50%)
        """
        # Device setup
        if device == 'auto':
            self.device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
        else:
            self.device = torch.device(device)
        
        print(f"Loading model on {self.device}...")
        
        # Load model
        self.model = create_ensemble(num_classes=len(CLASS_NAMES))
        checkpoint = torch.load(model_path, map_location=self.device, weights_only=False)
        
        self.model.efficientnet.load_state_dict(checkpoint['efficientnet_state_dict'])
        self.model.densenet.load_state_dict(checkpoint['densenet_state_dict'])
        
        self.model.to(self.device)
        self.model.eval()
        
        # Thresholds
        self.healthy_threshold = healthy_threshold
        self.uncertainty_threshold = uncertainty_threshold
        
        # Transform
        self.transform = get_transforms('test', img_size=224)
        
        print("✓ Model loaded successfully!")
        print(f"  EfficientNetB0 Recall: {checkpoint.get('efficientnet_recall', 'N/A')}")
        print(f"  DenseNet121 Recall: {checkpoint.get('densenet_recall', 'N/A')}")
        print(f"\nSafety Settings:")
        print(f"  Healthy Threshold: {self.healthy_threshold*100:.0f}%")
        print(f"  Uncertainty Threshold: {self.uncertainty_threshold*100:.0f}%")
    
    def preprocess_image(self, image: Union[str, Path, Image.Image]) -> torch.Tensor:
        """
        Load and preprocess an image.
        
        Args:
            image: Path to image or PIL Image
        
        Returns:
            Preprocessed tensor
        """
        if isinstance(image, (str, Path)):
            image = Image.open(image).convert('RGB')
        elif not isinstance(image, Image.Image):
            raise ValueError("Image must be a file path or PIL Image")
        
        # Apply transforms
        img_tensor = self.transform(image)
        
        # Add batch dimension
        img_tensor = img_tensor.unsqueeze(0)
        
        return img_tensor
    
    def predict(
        self,
        image: Union[str, Path, Image.Image],
        return_details: bool = False
    ) -> Union[Dict, str]:
        """
        Predict disease from fecal image with safety-first logic.
        
        Args:
            image: Image to classify (path or PIL Image)
            return_details: If True, return detailed results; else return simple classification
        
        Returns:
            If return_details=False: Classification string ('Healthy', 'Salmonella', 'ISOLATE', etc.)
            If return_details=True: Dictionary with full prediction details
        """
        # Preprocess
        img_tensor = self.preprocess_image(image).to(self.device)
        
        # Predict with safety vote
        with torch.no_grad():
            results = self.model.predict_with_safety_vote(
                img_tensor,
                healthy_threshold=self.healthy_threshold,
                uncertainty_threshold=self.uncertainty_threshold
            )
        
        # Extract results
        prediction_idx = results['prediction'][0].item()
        should_isolate = results['should_isolate'][0].item()
        
        efficient_probs = results['efficientnet_probs'][0].cpu().numpy()
        dense_probs = results['densenet_probs'][0].cpu().numpy()
        
        # Get confidence scores
        efficient_max_conf = results['confidence_scores']['efficientnet_max'][0].item()
        dense_max_conf = results['confidence_scores']['densenet_max'][0].item()
        efficient_healthy_conf = results['confidence_scores']['efficientnet_healthy'][0].item()
        dense_healthy_conf = results['confidence_scores']['densenet_healthy'][0].item()
        
        # Determine classification
        if prediction_idx == -1 or should_isolate:
            classification = "ISOLATE"
            risk_level = "HIGH RISK"
        else:
            classification = IDX_TO_CLASS[prediction_idx]
            risk_level = "LOW RISK" if classification == "Healthy" else "HIGH RISK"
        
        if not return_details:
            return classification
        
        # Determine action based on classification and isolation status
        if should_isolate:
            action = 'ISOLATE CHICKEN FOR INSPECTION'
        elif classification == "Healthy":
            action = 'Clear for flock'
        else:
            # Disease detected (Salmonella, Coccidiosis, New Castle Disease)
            action = 'ISOLATE CHICKEN - Disease Detected'
        
        # Detailed results
        return {
            'classification': classification,
            'risk_level': risk_level,
            'should_isolate': bool(should_isolate or classification != "Healthy"),  # Isolate if uncertain OR diseased
            'action': action,
            'models': {
                'efficientnet': {
                    'prediction': IDX_TO_CLASS[efficient_probs.argmax()],
                    'confidence': efficient_max_conf,
                    'healthy_confidence': efficient_healthy_conf,
                    'probabilities': {
                        CLASS_NAMES[i]: float(efficient_probs[i]) 
                        for i in range(len(CLASS_NAMES))
                    }
                },
                'densenet': {
                    'prediction': IDX_TO_CLASS[dense_probs.argmax()],
                    'confidence': dense_max_conf,
                    'healthy_confidence': dense_healthy_conf,
                    'probabilities': {
                        CLASS_NAMES[i]: float(dense_probs[i]) 
                        for i in range(len(CLASS_NAMES))
                    }
                }
            },
            'ensemble': {
                'agreement': IDX_TO_CLASS[efficient_probs.argmax()] == IDX_TO_CLASS[dense_probs.argmax()],
                'avg_confidence': (efficient_max_conf + dense_max_conf) / 2,
                'avg_healthy_confidence': (efficient_healthy_conf + dense_healthy_conf) / 2
            }
        }
    
    def predict_batch(
        self,
        images: list,
        return_details: bool = False
    ) -> list:
        """
        Predict on a batch of images.
        
        Args:
            images: List of image paths or PIL Images
            return_details: Return detailed results for each image
        
        Returns:
            List of predictions (one per image)
        """
        results = []
        for image in images:
            result = self.predict(image, return_details=return_details)
            results.append(result)
        return results
    
    def evaluate_safety(
        self,
        image: Union[str, Path, Image.Image]
    ) -> Tuple[bool, str]:
        """
        Quick safety evaluation.
        
        Args:
            image: Image to evaluate
        
        Returns:
            Tuple of (is_safe, reason)
        """
        result = self.predict(image, return_details=True)
        
        if result['should_isolate']:
            # Determine reason for isolation
            models_info = result['models']
            
            reasons = []
            
            # Check uncertainty
            if models_info['efficientnet']['confidence'] < self.uncertainty_threshold:
                reasons.append("EfficientNet uncertain")
            if models_info['densenet']['confidence'] < self.uncertainty_threshold:
                reasons.append("DenseNet uncertain")
            
            # Check healthy confidence
            if models_info['efficientnet']['healthy_confidence'] < self.healthy_threshold:
                reasons.append("EfficientNet low healthy confidence")
            if models_info['densenet']['healthy_confidence'] < self.healthy_threshold:
                reasons.append("DenseNet low healthy confidence")
            
            # Check disagreement
            if not result['ensemble']['agreement']:
                reasons.append("Models disagree")
            
            reason = "; ".join(reasons) if reasons else "Safety threshold not met"
            return False, reason
        
        return True, "Clear - High confidence in health status"


def main():
    """
    Example usage of the detector.
    """
    # Initialize detector
    detector = ChickenDiseaseDetector(
        model_path='outputs/ensemble_model.pth',
        healthy_threshold=0.80,
        uncertainty_threshold=0.50
    )
    
    # Example: Predict on a single image
    # Simple prediction
    result = detector.predict('path/to/fecal/image.jpg')
    print(f"Classification: {result}")
    
    # Detailed prediction
    detailed_result = detector.predict('path/to/fecal/image.jpg', return_details=True)
    print("\nDetailed Results:")
    print(f"Classification: {detailed_result['classification']}")
    print(f"Risk Level: {detailed_result['risk_level']}")
    print(f"Action: {detailed_result['action']}")
    print(f"\nEfficientNet:")
    print(f"  Prediction: {detailed_result['models']['efficientnet']['prediction']}")
    print(f"  Confidence: {detailed_result['models']['efficientnet']['confidence']:.2%}")
    print(f"\nDenseNet:")
    print(f"  Prediction: {detailed_result['models']['densenet']['prediction']}")
    print(f"  Confidence: {detailed_result['models']['densenet']['confidence']:.2%}")
    
    # Safety evaluation
    is_safe, reason = detector.evaluate_safety('path/to/fecal/image.jpg')
    print(f"\nSafety Check: {'SAFE' if is_safe else 'ISOLATE'}")
    print(f"Reason: {reason}")


if __name__ == '__main__':
    main()

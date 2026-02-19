"""
FastAPI service for Chicken Disease Detection Ensemble Model.
Provides REST API endpoints for disease prediction from fecal images.
"""

from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from PIL import Image
import io
import torch
from typing import Optional
import logging

from inference import ChickenDiseaseDetector

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Initialize FastAPI app
app = FastAPI(
    title="Tokkatot AI - Chicken Disease Detection",
    description="Safety-First Ensemble AI for detecting chicken diseases via fecal images",
    version="1.0.0"
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure this for production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Global detector instance
detector: Optional[ChickenDiseaseDetector] = None


class HealthResponse(BaseModel):
    """Health check response model"""
    status: str
    model_loaded: bool
    device: str


class PredictionResponse(BaseModel):
    """Prediction response model"""
    classification: str
    risk_level: str
    should_isolate: bool
    action: str
    confidence: float
    details: Optional[dict] = None


@app.on_event("startup")
async def startup_event():
    """Initialize the model on startup"""
    global detector
    try:
        logger.info("Loading ensemble model...")
        detector = ChickenDiseaseDetector(
            model_path='outputs/ensemble_model.pth',
            device='auto',
            healthy_threshold=0.80,
            uncertainty_threshold=0.50
        )
        logger.info("✓ Model loaded successfully!")
    except Exception as e:
        logger.error(f"Failed to load model: {str(e)}")
        raise


@app.get("/", response_model=dict)
async def root():
    """Root endpoint"""
    return {
        "message": "Tokkatot AI - Chicken Disease Detection API",
        "version": "1.0.0",
        "endpoints": {
            "health": "/health",
            "predict": "/predict",
            "predict_detailed": "/predict/detailed"
        }
    }


@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy" if detector is not None else "unhealthy",
        "model_loaded": detector is not None,
        "device": str(detector.device) if detector else "none"
    }


@app.post("/predict", response_model=PredictionResponse)
async def predict(file: UploadFile = File(...)):
    """
    Predict disease from fecal image (simple response).
    
    Args:
        file: Image file (JPEG, PNG)
    
    Returns:
        Simple classification result with action recommendation
    """
    if detector is None:
        raise HTTPException(status_code=503, detail="Model not loaded")
    
    # Validate file type
    if not file.content_type.startswith('image/'):
        raise HTTPException(status_code=400, detail="File must be an image")
    
    try:
        # Read image
        contents = await file.read()
        image = Image.open(io.BytesIO(contents)).convert('RGB')
        
        # Get prediction with details (we need some details even for simple response)
        result = detector.predict(image, return_details=True)
        
        # Return simplified response
        return {
            "classification": result['classification'],
            "risk_level": result['risk_level'],
            "should_isolate": result['should_isolate'],
            "action": result['action'],
            "confidence": result['ensemble']['avg_confidence']
        }
    
    except Exception as e:
        logger.error(f"Prediction error: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Prediction failed: {str(e)}")


@app.post("/predict/detailed", response_model=PredictionResponse)
async def predict_detailed(file: UploadFile = File(...)):
    """
    Predict disease from fecal image (detailed response).
    
    Args:
        file: Image file (JPEG, PNG)
    
    Returns:
        Detailed classification result with individual model predictions and probabilities
    """
    if detector is None:
        raise HTTPException(status_code=503, detail="Model not loaded")
    
    # Validate file type
    if not file.content_type.startswith('image/'):
        raise HTTPException(status_code=400, detail="File must be an image")
    
    try:
        # Read image
        contents = await file.read()
        image = Image.open(io.BytesIO(contents)).convert('RGB')
        
        # Get detailed prediction
        result = detector.predict(image, return_details=True)
        
        # Return full response with details
        return {
            "classification": result['classification'],
            "risk_level": result['risk_level'],
            "should_isolate": result['should_isolate'],
            "action": result['action'],
            "confidence": result['ensemble']['avg_confidence'],
            "details": result
        }
    
    except Exception as e:
        logger.error(f"Prediction error: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Prediction failed: {str(e)}")


@app.post("/evaluate/safety")
async def evaluate_safety(file: UploadFile = File(...)):
    """
    Quick safety evaluation for fecal image.
    
    Args:
        file: Image file (JPEG, PNG)
    
    Returns:
        Safety assessment with reason
    """
    if detector is None:
        raise HTTPException(status_code=503, detail="Model not loaded")
    
    # Validate file type
    if not file.content_type.startswith('image/'):
        raise HTTPException(status_code=400, detail="File must be an image")
    
    try:
        # Read image
        contents = await file.read()
        image = Image.open(io.BytesIO(contents)).convert('RGB')
        
        # Evaluate safety
        is_safe, reason = detector.evaluate_safety(image)
        
        return {
            "is_safe": is_safe,
            "status": "SAFE" if is_safe else "ISOLATE",
            "reason": reason
        }
    
    except Exception as e:
        logger.error(f"Safety evaluation error: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Evaluation failed: {str(e)}")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)

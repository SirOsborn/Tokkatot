# 🔬 AI Disease Detection Feature

This feature adds AI-powered chicken disease identification through feces analysis to your Tokkatot Smart Poultry System using a **PyTorch ensemble model** (99% combined accuracy) combining EfficientNetB0 and DenseNet121.

## 🎯 Overview

Farmers can use their mobile phones to:
- Take photos of chicken droppings
- Get instant AI-powered disease predictions with confidence scores
- Receive treatment recommendations based on diagnosis
- Access educational disease information
- Get high-confidence safety-first recommendations (ensemble voting)

## 🏗️ Architecture

```
Mobile Browser → Go API Gateway (6060) → FastAPI AI Service (8000) 
                                              ↓
                                    PyTorch Ensemble Models
                                    ├── EfficientNetB0 (98.05% recall)
                                    └── DenseNet121 (96.69% recall)
                                    ↓
                                    Ensemble Voting → 99% Accuracy
```

## 📋 Setup Instructions

### Prerequisites
- Python 3.12+ (recommended Python 3.12.1+)
- Go 1.19+ 
- PyTorch models in `ai-service/outputs/`:
  - `ensemble_model.pth` (ensemble combined model)
  - `checkpoints/EfficientNetB0_best.pth` (individual model)
  - `checkpoints/DenseNet121_best.pth` (individual model)

### Quick Start with Docker (Recommended)
```bash
cd ai-service
docker build -t tokkatot-ai-service .
docker run -p 8000:8000 tokkatot-ai-service
```

### Local Development Setup

Terminal 1 - AI Service:
```bash
cd ai-service

# Create virtual environment
python -m venv venv
source venv/bin/activate  # Linux/Mac
# OR
venv\Scripts\activate     # Windows

# Install dependencies
pip install -r requirements.txt

# Start FastAPI service
python app.py
# OR explicitly with Uvicorn:
uvicorn app:app --host 0.0.0.0 --port 8000 --reload
```

Terminal 2 - Go API Gateway:
```bash
cd middleware
go build -o tokkatot.exe .
./tokkatot.exe
```

## 📱 Usage

1. **Access the feature:** Visit `http://localhost:6060/disease-detection` or through mobile app
2. **Take/upload photo:** Use camera or upload existing image (PNG, JPEG)
3. **Get prediction:** AI ensemble analyzes the image (~1-3 seconds)
4. **Review results:** 
   - Disease classification with confidence scores from both models
   - Ensemble voting result (most confident prediction)
   - Treatment and prevention advice
5. **Follow recommendations:** Get actionable guidance

## 🔧 API Endpoints

### AI Service Endpoints (Port 8000)
- `GET /health` - Check AI service status and model loading
- `POST /predict` - Upload image for disease prediction (returns ensemble result)
- `POST /predict/detailed` - Upload image with detailed per-model confidence scores

### Integration with Go API Gateway (Port 6060)
- `GET /api/ai/health` - Forward to AI service health check
- `POST /api/predict` - Forward to AI prediction endpoint
- `POST /api/predict/detailed` - Forward to detailed prediction endpoint

See [IG_SPECIFICATIONS_AI_SERVICE.md](./IG_SPECIFICATIONS_AI_SERVICE.md) for detailed API specifications.

## 🎛️ Technical Specifications

### Model Architecture
- **Type:** PyTorch Ensemble (voting-based)
- **Primary Models:**
  - **EfficientNetB0:** 98.05% recall - Fast, efficient baseline
  - **DenseNet121:** 96.69% recall - Dense connections for nuanced features
- **Combined Accuracy:** 99% (via ensemble voting)
- **Framework:** PyTorch 2.0.0+
- **Input Size:** 224x224x3 (RGB)
- **Classes:** 4 disease categories
- **Python Version:** 3.12+

### Supported Disease Classes
This ensemble model detects:
- **Coccidiosis** - Parasitic infection affecting intestines
- **Healthy** - Normal, healthy droppings
- **Newcastle Disease** - Viral respiratory and nervous system disease  
- **Salmonella** - Bacterial infection

### Performance Specifications
- **Prediction Time:** ~1-3 seconds on CPU (GPU: <500ms)
- **Memory Usage:** ~800-1200MB (PyTorch + ensemble models)
- **Input Processing:** Automatic resize to 224x224, normalization
- **Model Size:** 47.2 MB (ensemble), individual checkpoints ~100MB each
- **Inference Precision:** Float32 (FP32) for maximum accuracy

## 🎛️ Raspberry Pi Optimization

### Resource Management (Tested on Pi 4, 4GB RAM)
- **CPU Usage:** ~30-50% during ensemble inference
- **Memory:** ~800-1200MB for PyTorch + ensemble models
- **Storage:** ~200-250MB for model files (checkpoints + ensemble)
- **Inference Time:** ~5-10 seconds (Pi 4 CPU only)

### Performance Tips

1. **Use GPU Acceleration (if available):**
```bash
# Install PyTorch with CUDA support
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118

# Or for CPU-only (smaller download):
pip install torch torchvision --index-url https://download.pytorch.org/whl/cpu
```

2. **Model Quantization (future optimization):**
```python
# Convert to quantized model for ~4x faster inference
import torch
quantized_model = torch.quantization.quantize_dynamic(
    model, {torch.nn.Linear}, dtype=torch.qint8
)
torch.save(quantized_model.state_dict(), 'outputs/ensemble_model_quantized.pth')
```

3. **Memory Optimization in app.py:**
```python
import torch
import gc

# Clear cache between predictions
def cleanup():
    gc.collect()
    torch.cuda.empty_cache()  # if using GPU

# Call after each prediction
cleanup()
```

4. **Image Preprocessing on Client:**
Resize images before upload to reduce bandwidth:
```python
from PIL import Image
img = Image.open('photo.jpg')
img.thumbnail((224, 224))  # Pre-resize before upload
img.save('photo_small.jpg')
```

## 🔒 Security Considerations

- File upload validation (size max 5MB, types: PNG/JPEG only)
- User authentication required via Go API Gateway
- Rate limiting: 100 requests per minute per user (Go middleware)
- Secure model file storage in `outputs/` (not exposed publicly)
- Input validation: Verify image dimensions before processing
- Error handling: No sensitive model paths exposed in error messages

## 🛠️ Troubleshooting

### Common Issues

**AI Service not starting:**
- Ensure Python 3.12+ is installed: `python --version`
- Check virtual environment: `venv/` exists
- Verify dependencies installed: `pip list | grep torch`
- Check that FastAPI/Uvicorn are installed: `pip show fastapi uvicorn`
- Ensure model files exist: `ls outputs/ensemble_model.pth`
- Check port 8000 not in use: `netstat -an | grep 8000`

**Model loading errors:**
- Verify PyTorch 2.0+ installation: `pip show torch`
- Check model files in `outputs/`:
  - `ensemble_model.pth` (required)
  - `checkpoints/EfficientNetB0_best.pth` (required)
  - `checkpoints/DenseNet121_best.pth` (required)
- Verify file permissions: `chmod 644 outputs/*.pth`
- Test model loading manually:
```python
import torch
from models import EnsembleModel
model = EnsembleModel()
model.load_state_dict(torch.load('outputs/ensemble_model.pth'))
print("Model loaded successfully")
```

**Low accuracy predictions:**
- Verify image quality (good lighting, clear focus on droppings)
- Ensure image is RGB (not grayscale or RGBA)
- Check image size is at least 224x224 (or it will be resized)
- Validate image preprocessing matches training pipeline
- Verify you're using the correct model files

**Import errors for PyTorch:**
- Reinstall PyTorch: `pip install --upgrade torch torchvision`
- Check CUDA compatibility if using GPU: `python -c "import torch; print(torch.cuda.is_available())"`
- For CPU-only: Reinstall PyTorch CPU version

**High memory usage:**
- Monitor memory: `ps aux | grep python` (Linux) or Task Manager (Windows)
- Enable garbage collection: Add `gc.collect()` after predictions
- Reduce batch size (if batch processing)
- Use CPU instead of GPU if GPU memory exhausted

**Slow predictions:**
- Check if running on CPU vs GPU: `torch.cuda.is_available()`
- Optimize image preprocessing (resize before upload)
- Verify no other CPU-intensive processes running
- Consider using quantized model for Pi deployment
- Check disk I/O: Ensure SSD, not spinning disk

### Health Checks

Monitor service health:
```bash
# Check FastAPI service health
curl http://127.0.0.1:8000/health

# Test prediction endpoint
curl -X POST -F "image=@test_image.jpg" http://127.0.0.1:8000/predict

# Check Go API Gateway
curl http://127.0.0.1:6060/api/ai/health
```

### Debug Mode

Enable debug logging in `app.py`:
```python
import logging
logging.basicConfig(level=logging.DEBUG)

# Start with Uvicorn debug:
uvicorn app:app --host 0.0.0.0 --port 8000 --reload
```

## 🔄 Model Updates

To update the AI ensemble model:
1. Replace files in `ai-service/outputs/`:
   - `ensemble_model.pth` (main ensemble weights)
   - `checkpoints/EfficientNetB0_best.pth` (EfficientNetB0 checkpoint)
   - `checkpoints/DenseNet121_best.pth` (DenseNet121 checkpoint)
2. Restart the AI service:
```bash
# If running with Docker
docker restart tokkatot-ai-service

# If running locally
# Kill current process and run: python app.py
```

No changes needed to Go API Gateway - it will automatically use new model files.

### Model Versioning Strategy
Consider versioning the models:
```
ai-service/
├── outputs/
│   ├── v1.0/
│   │   ├── ensemble_model.pth
│   │   └── checkpoints/
│   ├── v2.0/
│   │   ├── ensemble_model.pth
│   │   └── checkpoints/
│   └── current/ (symlink to latest version)
├── app.py (references outputs/current/)
```

## 📈 Monitoring & Analytics

### System Metrics
- Prediction response time (target: <3 seconds on CPU)
- Model ensemble confidence scores (both EfficientNetB0 and DenseNet121)
- Prediction consistency (agreement between models)
- Error rates and types (invalid images, timeouts)
- Resource usage (CPU %, memory %, disk I/O)
- Model agreement percentage (when both models agree vs disagree)
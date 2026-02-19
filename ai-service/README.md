# Tokkatot AI Service

**Safety-First Ensemble AI for Chicken Disease Detection**

## 📁 Directory Structure

```
ai-service/
├── app.py                          # FastAPI service entry point
├── inference.py                    # Inference module
├── models.py                       # EfficientNetB0 & DenseNet121 architectures
├── data_utils.py                   # Data transforms and utilities
├── pyproject.toml                  # Python dependencies
├── Dockerfile                      # Docker container configuration
├── docker-compose.yml              # Docker Compose orchestration
├── .dockerignore                   # Docker build exclusions
├── DOCKER_DEPLOYMENT.md            # Detailed deployment guide
├── MODEL_CARD_ENSEMBLE.md          # Ensemble model documentation
├── MODEL_CARD_EfficientNetB0.md    # EfficientNetB0 model card
├── MODEL_CARD_DenseNet121.md       # DenseNet121 model card
└── outputs/
    ├── ensemble_model.pth          # Ensemble model weights (47.2 MB)
    └── checkpoints/
        ├── DenseNet121_best.pth    # DenseNet121 weights
        └── EfficientNetB0_best.pth # EfficientNetB0 weights
```

## 🚀 Quick Start

### Option 1: Docker (Recommended for Production)

```bash
# Build and run with Docker Compose
cd C:\Users\PureGoat\tokkatot\ai-service
docker-compose up -d

# Access the API
# - API Documentation: http://localhost:8000/docs
# - Health Check: http://localhost:8000/health
```

### Option 2: Local Development

```bash
# Create virtual environment
python -m venv .venv
.venv\Scripts\Activate.ps1

# Install dependencies
pip install -e .

# Run the service
python -m uvicorn app:app --host 0.0.0.0 --port 8000
```

## 📊 API Endpoints

### Health Check
```bash
curl http://localhost:8000/health
```

### Simple Prediction
```bash
curl -X POST "http://localhost:8000/predict" \
  -H "Content-Type: multipart/form-data" \
  -F "file=@path/to/image.jpg"
```

### Detailed Prediction
```bash
curl -X POST "http://localhost:8000/predict/detailed" \
  -H "Content-Type: multipart/form-data" \
  -F "file=@path/to/image.jpg"
```

### Safety Evaluation
```bash
curl -X POST "http://localhost:8000/evaluate/safety" \
  -H "Content-Type: multipart/form-data" \
  -F "file=@path/to/image.jpg"
```

## 🎯 Model Performance

- **Overall Accuracy:** 99%
- **Overall Recall:** 99%
- **EfficientNetB0 Recall:** 98.05%
- **DenseNet121 Recall:** 96.69%

### Safety-First Decision Logic

The system **isolates** chickens if ANY condition is met:

1. **Uncertainty:** Either model's max confidence < 50%
2. **Safety Vote:** Either model's healthy confidence < 80%
3. **Disagreement:** Models disagree AND either predicts disease

## 🦠 Detected Diseases

- **Healthy** - Normal fecal matter
- **Salmonella** - Bacterial infection
- **Coccidiosis** - Parasitic infection
- **New Castle Disease** - Viral infection

## 🔗 Integration with Tokkatot Ecosystem

This AI service is designed to integrate seamlessly with:

- **Embedded System** (`embedded/`) - Raspberry Pi image capture
- **Middleware** (`middleware/`) - IoT message routing
- **Frontend** (`frontend/`) - User dashboard

### Example Integration

```python
import requests

# From embedded device or middleware
with open("captured_fecal_image.jpg", "rb") as f:
    files = {"file": f}
    response = requests.post(
        "http://ai-service:8000/predict/detailed",
        files=files
    )
    result = response.json()
    
    if result["should_isolate"]:
        # Trigger isolation mechanism
        isolate_chicken(result["classification"])
```

## 📦 Dependencies

Core dependencies (from pyproject.toml):
- PyTorch >= 2.0.0
- TorchVision >= 0.15.0
- FastAPI
- Uvicorn
- Pillow >= 10.0.0
- NumPy >= 1.24.0

## 🐳 Docker Deployment

See [DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md) for:
- Cloud deployment (AWS ECS, Google Cloud Run, Azure)
- Resource requirements
- Performance monitoring
- Production configuration

## 📝 Model Documentation

- [Ensemble Model Card](MODEL_CARD_ENSEMBLE.md) - Combined system
- [EfficientNetB0 Card](MODEL_CARD_EfficientNetB0.md) - Fast, edge-optimized
- [DenseNet121 Card](MODEL_CARD_DenseNet121.md) - Robust, feature-rich

## ⚡ Performance

Based on Docker CPU deployment testing:
- **Inference Time:** 150-330ms per image (avg 244ms)
- **Memory Usage:** ~648 MB
- **CPU Usage:** < 1% (idle), brief spikes during inference

## 🔐 Security Notes

- FastAPI runs as non-root user in Docker
- CORS configured (update for production)
- Health checks enabled
- Input validation on image uploads

## 📄 License

**© 2026 Tokkatot. All Rights Reserved.**

---

**Developed exclusively for Tokkatot Smart Chicken Farming Solutions**

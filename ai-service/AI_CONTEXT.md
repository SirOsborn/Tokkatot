# 🤖 AI Context: AI/Disease Detection Service

**Directory**: `ai-service/`  
**Your Role**: AI models, disease prediction, FastAPI REST API  
**Tech Stack**: Python 3.12, PyTorch 2.0, FastAPI, Uvicorn  

---

## 🎯 What You're Building

**Ensemble AI Model for Chicken Disease Detection**
- **Input**: PNG/JPEG images of chicken feces (~5MB max)
- **Output**: Disease classification + confidence scores + treatment recommendations
- **Accuracy**: 99% (via ensemble voting: EfficientNetB0 + DenseNet121)
- **Performance**: 1-3 seconds on CPU, <500ms on GPU

**API Endpoints** (FastAPI, Port 8000):
1. `GET /health` - Service health & model loading status
2. `POST /predict` - Disease prediction (simple response)
3. `POST /predict/detailed` - Detailed per-model confidence scores

---

## 📁 File Structure

```
ai-service/
├── app.py                 # FastAPI server, endpoint definitions
├── models.py             # PyTorch model architectures (EfficientNetB0, DenseNet121)
├── inference.py          # Ensemble model loading & prediction logic
├── data_utils.py         # Image preprocessing, transforms, class names
├── requirements.txt      # Python dependencies (FastAPI, torch, etc)
├── pyproject.toml        # Python project config
├── Dockerfile            # Docker build (Python 3.12-slim, model copying)
├── docker-compose.yml    # Docker Compose config
├── AI_FEATURE_README.md  # User-facing documentation
├── AI_CONTEXT.md         # This file
├── outputs/              # Model files (DO NOT COMMIT - use .gitignore)
│   ├── ensemble_model.pth           # Ensemble weights
│   └── checkpoints/
│       ├── EfficientNetB0_best.pth
│       └── DenseNet121_best.pth
└── .gitignore           # Exclude model files from git
```

---

## 🚀 Getting Started

### Local Development

```bash
cd ai-service

# Create virtual environment
python -m venv venv
source venv/bin/activate  # Linux/Mac
# OR: venv\Scripts\activate  # Windows

# Install dependencies
pip install -r requirements.txt

# Start FastAPI service
python app.py
# OR: uvicorn app:app --host 0.0.0.0 --port 8000 --reload

# Test health
curl http://localhost:8000/health

# Test prediction
curl -X POST -F "image=@sample_image.jpg" http://localhost:8000/predict
```

### Docker

```bash
# Build image
docker build -t tokkatot-ai .

# Run container
docker run -p 8000:8000 tokkatot-ai

# Or use docker-compose
docker-compose up
```

---

## 📊 Model Architecture

### Ensemble Structure
- **Primary**: EfficientNetB0 (98.05% recall)
  - Fast, lightweight, good for edge devices
- **Secondary**: DenseNet121 (96.69% recall)
  - Dense connections capture nuanced features
- **Voting**: Select prediction with highest combined confidence (99% accuracy)

### Inference Flow
```
User uploads image
    ↓
Validate: PNG/JPEG, max 5MB
    ↓
Preprocess: Resize to 224x224, normalize
    ↓
Forward pass (parallel):
   ├─ EfficientNetB0 → confidence score
   └─ DenseNet121 → confidence score
    ↓
Ensemble voting: Select highest confidence
    ↓
Safety check: If confidence < 50%, return "uncertain"
    ↓
Return: disease, confidence, treatment_options
```

### Supported Diseases
1. **Coccidiosis** - Parasitic infection affecting intestines
2. **Healthy** - Normal, healthy droppings
3. **Newcastle Disease** - Viral respiratory and nervous system disease
4. **Salmonella** - Bacterial infection

---

## 🔧 Key Functions

### `app.py`
- `startup_event()` - Initialize model on service startup
- `health_check()` - Return model status & device info
- `predict()` - Accept image, return disease prediction
- `predict_detailed()` - Accept image, return detailed per-model scores

### `inference.py` - `ChickenDiseaseDetector` class
- `__init__()` - Load ensemble model from `.pth` file
- `preprocess_image()` - PIL Image → PyTorch tensor
- `predict()` - Inference on single image
- `evaluate_safety()` - Safety-first decision logic

### `models.py`
- `create_ensemble()` - Instantiate ensemble model
- `EfficientNetB0Wrapper` - Model architecture
- `DenseNet121Wrapper` - Model architecture

### `data_utils.py`
- `CLASS_NAMES` - List of 4 disease classes
- `get_transforms()` - Preprocessing & augmentation
- `validate_image()` - Check image format/size

---

## 📝 Code Guidelines

### ✅ DO:
- Use type hints on all functions
- Add docstrings explaining inference logic
- Validate image uploads (size, format, dimensions)
- Return structured JSON responses
- Include confidence scores in responses
- Handle exceptions gracefully (return 400/500 with clear message)
- Log important events (model loading, predictions, errors)
- Use `.env` for any configuration

### ❌ DON'T:
- Hardcode model paths (use `'outputs/ensemble_model.pth'` as relative path)
- Expose internal model details in error messages
- Accept images larger than 5MB
- Trust file extensions (check magic bytes for PNG/JPEG)
- Return None or empty responses
- Commit `.pth` model files to git
- Hardcode API URLs or config values

---

## 🔒 Security Checklist

- ✅ Validate file uploads (size max 5MB, type PNG/JPEG)
- ✅ No exception details in error responses
- ✅ No model file paths in responses
- ✅ Rate limiting on `/predict` endpoint (100 req/min per user)
- ✅ JWT token validation (via Go API Gateway)
- ✅ Input preprocessing prevents adversarial attacks (normalization)
- ✅ Model files never pushed to git (`.gitignore`)
- ✅ Secrets in `.env` (not hardcoded)

---

## 📊 Performance Targets

| Metric | Target | Actual |
|--------|--------|--------|
| Inference Time (CPU) | < 3 seconds | ~1.2s |
| Inference Time (GPU) | < 500ms | ~300ms |
| Model Load Time | < 60 seconds | ~45s |
| Memory Usage | < 2GB | ~1.2GB |
| Response Time (API) | < 100ms | ~50ms |
| Model Size | < 100MB | 47.2MB |

---

## 🧪 Testing

### Unit Tests
```python
def test_image_preprocessing():
    # Test tensor output shape
    tensor = detector.preprocess_image('test_image.jpg')
    assert tensor.shape == (1, 3, 224, 224)

def test_prediction_response_structure():
    # Test JSON response has required fields
    result = detector.predict(image)
    assert 'disease' in result
    assert 'confidence' in result
    assert 'recommendation' in result
```

### Integration Tests
```bash
curl -X POST -F "image=@healthy.jpg" http://localhost:8000/predict
# Expected: {"disease": "Healthy", "confidence": 0.99, ...}

curl -X POST -F "image=@coccidiosis.jpg" http://localhost:8000/predict
# Expected: {"disease": "Coccidiosis", "confidence": 0.98, ...}
```

---

## 📈 What You Can Change

✅ **You have authority over**:
- PyTorch model implementation
- API endpoint logic
- Docker configuration
- Error handling & validation
- Performance optimizations
- Testing & documentation for AI service

❓ **You should ask before changing**:
- Model architecture (should discuss with ML team)
- API response format (impacts Go API & frontend)
- Supported diseases (impacts database schema)
- Model training/retraining (coordinate with data team)

---

## 🔗 Integration Points

**With Go API Gateway** (Port 6060):
```
Go API (Port 6060)
    ↓
[HTTP GET /api/ai/health]
    ↓
FastAPI Health Check (Port 8000)
    ↓
Returns: {"status": "healthy", "model_loaded": true}
```

**With Go API** (for predictions):
```
Go API receives user image upload
    ↓
[HTTP POST /api/ai/predict with multipart image]
    ↓
FastAPI Prediction Logic
    ↓
[Returns JSON with disease, confidence, treatment]
    ↓
Go API stores result in PostgreSQL & broadcasts via WebSocket
```

**Database Integration**:
- Go API stores predictions in `prediction_logs` table
- Records: `user_id`, `device_id`, `prediction_id`, `disease`, `confidence`, `timestamp`
- Can query predictions via `GET /farms/{farm_id}/predictions`

---

## 🆘 Common Issues & Solutions

### Issue: Model not loading
```
Error: torch.load() failed
```
**Fix**: Check that model files exist in `outputs/` and file paths are correct in `inference.py`

### Issue: Out of memory
```
Error: CUDA out of memory
```
**Fix**: Use CPU instead (`device='cpu'`), or reduce batch size, or optimize model quantization

### Issue: Image preprocessing fails
```
Error: PIL cannot open image
```
**Fix**: Validate image format and convert to RGB before processing

### Issue: Predictions are inconsistent
```
Different results for same image
```
**Fix**: Ensure model is in `.eval()` mode, disable dropout & batch norm during inference

---

## 📚 Key Documents

- `IG_SPECIFICATIONS_AI_SERVICE.md` - Full AI service specification (endpoints, architecture)
- `AI_FEATURE_README.md` - User guide (setup, usage, troubleshooting)
- `01_SPECIFICATIONS_ARCHITECTURE.md` - System-wide architecture (how AI fits in)

---

## 🎯 Your Next Tasks

1. **Implement** according to `IG_SPECIFICATIONS_AI_SERVICE.md`
2. **Test locally** - Verify health endpoint, test predictions
3. **Build Docker image** - Ensure model files copy correctly
4. **Document** - Keep API_FEATURE_README.md current
5. **Integrate** - Coordinate with Go API for endpoint routing
6. **Monitor** - Set up logging for prediction errors/anomalies

---

**Happy coding! 🚀 If unexpected issues arise, check the spec first, then ask the team.**

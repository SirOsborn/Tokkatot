# 🤖 AI Context: PyTorch Disease Detection Service

**Component**: `ai-service/` - AI-powered chicken disease detection  
**Tech Stack**: Python 3.12, PyTorch 2.0, FastAPI, Uvicorn  
**Purpose**: Ensemble model inference (EfficientNetB0 + DenseNet121) for disease classification from feces images  

---

## 📖 Read First

**Before reading this file**, understand the project context:
- **Project overview**: Read [`../AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) for business model, farmer needs, why AI disease detection matters
- **Full AI service spec**: See [`../docs/implementation/AI_SERVICE.md`](../docs/implementation/AI_SERVICE.md) for complete API specs, model details, deployment

**This file contains**: PyTorch-specific patterns, model loading, FastAPI endpoints, Docker deployment

---

## 📚 Full Documentation

| Document | Purpose |
|----------|---------|
| [`docs/implementation/AI_SERVICE.md`](../docs/implementation/AI_SERVICE.md) | Complete AI service specs (endpoints, model architecture, training pipeline) |
| [`docs/implementation/API.md`](../docs/implementation/API.md) | How Go middleware calls `/predict` endpoint |
| [`ai-service/AI_FEATURE_README.md`](./AI_FEATURE_README.md) | User-facing guide (setup, testing, troubleshooting) |
| [`AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) | Why farmers need disease detection (prevent spread, early treatment) |

---

## 📁 Quick File Reference

```
ai-service/
├── app.py                      # FastAPI server (3 endpoints: /health, /predict, /predict/detailed)
├── inference.py                # ChickenDiseaseDetector class (model loading, ensemble voting)
├── models.py                   # PyTorch architectures (EfficientNetB0Wrapper, DenseNet121Wrapper)
├── data_utils.py               # Image preprocessing, class definitions, transforms
├── requirements.txt            # PyTorch 2.0, FastAPI, Uvicorn, Pillow, Pydantic
├── pyproject.toml              # Project metadata, dependencies
├── Dockerfile                  # Python 3.12-slim, COPY outputs/, health checks
├── docker-compose.yml          # Port 8000, 2 CPU / 4GB RAM limits
├── .env                        # Config (NOT COMMITTED - see .gitignore)
├── outputs/
│   ├── ensemble_model.pth      # 47.2 MB ensemble weights (PROPRIETARY - NOT COMMITTED)
│   └── checkpoints/
│       ├── EfficientNetB0_best.pth
│       └── DenseNet121_best.pth
└── AI_CONTEXT.md               # This file
```

**CRITICAL**: `outputs/*.pth` files are **NOT in Git** (`.gitignore`'d). Model files must exist locally for app to start.

---

## 🎯 AI Service Purpose (Farmer Context)

**Why This Service Exists**:
- Cambodian farmers cannot afford veterinarians for every sick chicken
- Early disease detection prevents outbreaks across entire flock (100+ birds)
- Visual inspection of feces is traditional but unreliable (farmers not trained in pathology)
- AI model provides 99% accurate diagnosis in <3 seconds

**Farmer Workflow**:
1. Farmer notices sick chicken (lethargy, diarrhea)
2. Takes photo of feces with phone camera
3. Uploads to Tokkatot app (frontend calls Go API → calls this service)
4. Gets disease name + confidence + treatment steps in Khmer
5. Follows recommendations (isolate bird, administer treatment)

**Impact**: Early detection saves 80% of flocks from preventable outbreaks (see [`../docs/02_SPECIFICATIONS_REQUIREMENTS.md`](../docs/02_SPECIFICATIONS_REQUIREMENTS.md) - Business Goals).

---

## 🔧 Model Architecture (Ensemble)

**Structure**: 2 PyTorch models vote on final prediction

### Model 1: EfficientNetB0
- **Architecture**: Compound scaling (width, depth, resolution)
- **Input**: 224x224 RGB image
- **Output**: 4-class probabilities (Coccidiosis, Healthy, Newcastle, Salmonella)
- **Recall**: 98.05% (excellent at detecting diseases)
- **File**: `outputs/checkpoints/EfficientNetB0_best.pth`
- **Why**: Lightweight, fast inference (good for edge devices)

### Model 2: DenseNet121
- **Architecture**: Dense connections between layers (feature reuse)
- **Input**: 224x224 RGB image
- **Output**: 4-class probabilities
- **Recall**: 96.69%
- **File**: `outputs/checkpoints/DenseNet121_best.pth`
- **Why**: Captures nuanced features (useful for similar-looking diseases)

### Ensemble Voting Logic
```python
def ensemble_predict(image):
    # 1. Run both models in parallel
    efficientnet_probs = model1.forward(image)  # [0.02, 0.95, 0.01, 0.02] (4 classes)
    densenet_probs = model2.forward(image)       # [0.05, 0.90, 0.03, 0.02]
    
    # 2. Average probabilities
    avg_probs = (efficientnet_probs + densenet_probs) / 2
    
    # 3. Select highest confidence class
    predicted_class = argmax(avg_probs)
    confidence = max(avg_probs)
    
    # 4. Safety check
    if confidence < 0.50:
        return "uncertain", confidence, "Please retake photo in better lighting"
    
    return CLASS_NAMES[predicted_class], confidence, get_treatment(predicted_class)
```

**Result**: 99% aggregate accuracy (better than individual models).

---

## 🌐 FastAPI Endpoints

### 1. `GET /health`
**Purpose**: Service health check (used by Go middleware & Docker)

**Response**:
```json
{
  "status": "healthy",
  "model_loaded": true,
  "device": "cpu",
  "timestamp": "2025-02-01T12:00:00Z"
}
```

**Implementation** (`app.py`):
```python
@app.get("/health")
async def health_check():
    return {
        "status": "healthy" if detector.model_loaded else "unavailable",
        "model_loaded": detector.model_loaded,
        "device": str(detector.device),
        "timestamp": datetime.now().isoformat()
    }
```

### 2. `POST /predict`
**Purpose**: Disease prediction (main endpoint)

**Request**: Multipart form data with `image` field (PNG/JPEG, max 5MB)

**Response**:
```json
{
  "disease": "Coccidiosis",
  "confidence": 0.99,
  "recommendation": "Isolate affected birds immediately. Administer Amprolium 20% solution...",
  "treatment_options": ["Amprolium 20%", "Sulfonamides"]
}
```

**Implementation** (`app.py`):
```python
@app.post("/predict")
async def predict(image: UploadFile = File(...)):
    # 1. Validate image (size, format)
    if image.size > 5 * 1024 * 1024:  # 5MB limit
        raise HTTPException(400, "Image too large")
    
    # 2. Load image bytes
    image_bytes = await image.read()
    pil_image = Image.open(io.BytesIO(image_bytes)).convert('RGB')
    
    # 3. Run inference
    result = detector.predict(pil_image)
    
    # 4. Return structured response
    return {
        "disease": result['disease'],
        "confidence": float(result['confidence']),
        "recommendation": get_recommendation(result['disease']),
        "treatment_options": get_treatments(result['disease'])
    }
```

**Called By**: Go middleware endpoint `POST /api/ai/predict` (see [`../middleware/api/disease-detection.go`](../middleware/api))

### 3. `POST /predict/detailed`
**Purpose**: Detailed per-model confidence scores (for debugging, research)

**Response**:
```json
{
  "disease": "Coccidiosis",
  "confidence": 0.99,
  "model_predictions": {
    "EfficientNetB0": {
      "disease": "Coccidiosis",
      "confidence": 0.98
    },
    "DenseNet121": {
      "disease": "Coccidiosis",
      "confidence": 0.99
    }
  },
  "recommendation": "...",
  "treatment_options": [...]
}
```

**Use Case**: Research, model evaluation, debugging inconsistent predictions.

---

## 🐳 Docker Deployment

**Dockerfile**:
```dockerfile
FROM python:3.12-slim

WORKDIR /app

# Copy dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy model files (MUST exist during build)
COPY outputs/ outputs/

# Copy application code
COPY . .

# Health check (FastAPI /health endpoint)
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD curl -f http://localhost:8000/health || exit 1

# Start FastAPI server
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]
```

**Build & Run**:
```bash
cd ai-service

# Build image
docker build -t tokkatot-ai:latest .

# Run container
docker run -d -p 8000:8000 --name tokkatot-ai tokkatot-ai:latest

# Or use docker-compose
docker-compose up -d tokkatot-ai
```

**Resource Limits** (`docker-compose.yml`):
- CPU: 2 cores
- Memory: 4GB
- Restart policy: always (auto-restart on failure)

**Health Check**:
- Interval: 30 seconds
- Timeout: 10 seconds
- Retries: 3 before marking unhealthy

**Integration with Go Middleware**:
- Go middleware calls `http://localhost:8000/predict` (docker network or localhost)
- Timeout: 3 seconds CPU, 500ms GPU target
- Retry: 3 attempts with exponential backoff

---

## 🔧 Common PyTorch Patterns

### Model Loading (`inference.py`)

```python
class ChickenDiseaseDetector:
    def __init__(self, model_path='outputs/ensemble_model.pth', device='cpu'):
        self.device = torch.device(device)
        
        # Load ensemble model from disk
        checkpoint = torch.load(model_path, map_location=self.device)
        
        # Initialize architecture
        self.model = create_ensemble()
        self.model.load_state_dict(checkpoint['model_state_dict'])
        
        # Set to evaluation mode (disable dropout, batch norm)
        self.model.eval()
        self.model_loaded = True
    
    def preprocess_image(self, pil_image):
        # Resize to 224x224, normalize to ImageNet stats
        transform = transforms.Compose([
            transforms.Resize((224, 224)),
            transforms.ToTensor(),
            transforms.Normalize(mean=[0.485, 0.456, 0.406], 
                                 std=[0.229, 0.224, 0.225])
        ])
        return transform(pil_image).unsqueeze(0)  # Add batch dimension
    
    def predict(self, pil_image):
        # Disable gradient computation (inference only)
        with torch.no_grad():
            tensor = self.preprocess_image(pil_image).to(self.device)
            outputs = self.model(tensor)
            
            # Get probabilities from logits
            probs = torch.softmax(outputs, dim=1)
            
            # Best prediction
            confidence, predicted_class = torch.max(probs, dim=1)
            
            # Safety check
            if confidence < 0.50:
                return {
                    'disease': 'uncertain',
                    'confidence': confidence.item(),
                    'recommendation': 'Please retake photo in better lighting'
                }
            
            return {
                'disease': CLASS_NAMES[predicted_class.item()],
                'confidence': confidence.item(),
                'recommendation': get_recommendation(predicted_class.item())
            }
```

### Error Handling

```python
@app.post("/predict")
async def predict(image: UploadFile = File(...)):
    try:
        # Validate image
        if image.size > 5 * 1024 * 1024:
            raise HTTPException(status_code=400, detail="Image too large (max 5MB)")
        
        if image.content_type not in ["image/jpeg", "image/png"]:
            raise HTTPException(status_code=400, detail="Invalid image format (use PNG or JPEG)")
        
        # Load and predict
        image_bytes = await image.read()
        pil_image = Image.open(io.BytesIO(image_bytes)).convert('RGB')
        result = detector.predict(pil_image)
        
        return result
    
    except Exception as e:
        # Log error (don't expose internal details to user)
        logger.error(f"Prediction failed: {str(e)}")
        raise HTTPException(status_code=500, detail="Prediction failed")
```

---

## 🔒 Security Best Practices

- ✅ **File Upload Validation**: Max 5MB, PNG/JPEG only, check magic bytes (not just extension)
- ✅ **No Model Paths in Errors**: Never expose `outputs/ensemble_model.pth` in error messages
- ✅ **Rate Limiting**: FastAPI middleware limits `/predict` to 100 req/min per IP (prevent abuse)
- ✅ **No Secrets in Code**: Use `.env` for config (loaded via `python-dotenv`)
- ✅ **Model Files Not Committed**: `.gitignore` excludes `outputs/*.pth` (proprietary, 47MB+)
- ✅ **Input Normalization**: Prevents adversarial attacks (ImageNet stats normalization)

---

## 🆘 Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| **Model not loading** | Verify `outputs/ensemble_model.pth` exists, check file path in `inference.py` |
| **Out of memory (CUDA)** | Switch to CPU (`device='cpu'`), or reduce batch size, or use model quantization |
| **PIL cannot open image** | Validate image format before processing, convert to RGB mode |
| **Inconsistent predictions** | Ensure model is in `.eval()` mode, disable dropout/batch norm |
| **Slow inference (>3s)** | Use GPU, optimize model with TorchScript, reduce image size |
| **Docker build fails** | Ensure model files exist in `outputs/` before `docker build` |

---

## 🧪 Development Tasks

### Test Locally (Python venv)

```bash
cd ai-service

# Create virtual environment (REQUIRED)
python3 -m venv env
source env/bin/activate  # Windows: env\Scripts\activate

# Install dependencies
pip install --upgrade pip setuptools wheel
pip install -r requirements.txt

# Verify PyTorch installation
python3 -c "import torch; print(f'PyTorch {torch.__version__}')"

# Start FastAPI server
python3 app.py  # Runs on http://localhost:8000

# Test health endpoint
curl http://localhost:8000/health

# Test prediction with sample image
curl -X POST -F "image=@sample_healthy.jpg" http://localhost:8000/predict
```

### Test with Docker

```bash
# Build image (ensure model files exist in outputs/)
docker build -t tokkatot-ai:latest .

# Run container
docker run -d -p 8000:8000 --name tokkatot-ai tokkatot-ai:latest

# Check logs
docker logs tokkatot-ai

# Test endpoint
curl http://localhost:8000/health
```

### Add a New Disease Class

1. **Update `data_utils.py`**:
   ```python
   CLASS_NAMES = ['Coccidiosis', 'Healthy', 'Newcastle Disease', 'Salmonella', 'New Disease']
   ```

2. **Retrain models** with new dataset (coordinate with ML team)

3. **Update `inference.py`** treatment recommendations:
   ```python
   def get_recommendation(disease):
       if disease == 'New Disease':
           return "Specific treatment for new disease..."
   ```

4. **Update Go middleware** database schema (see [`../docs/implementation/DATABASE.md`](../docs/implementation/DATABASE.md) - `prediction_logs` table)

---

## 📘 Documentation Map

**AI Context Files** (component-specific guides):
- **This file**: [`ai-service/AI_CONTEXT.md`](./AI_CONTEXT.md) - PyTorch patterns, FastAPI endpoints
- [`middleware/AI_CONTEXT.md`](../middleware/AI_CONTEXT.md) - Go API, how it calls this service
- [`frontend/AI_CONTEXT.md`](../frontend/AI_CONTEXT.md) - Vue.js UI for disease detection page
- [`embedded/AI_CONTEXT.md`](../embedded/AI_CONTEXT.md) - ESP32 firmware (no direct connection to AI service)
- [`docs/AI_CONTEXT.md`](../docs/AI_CONTEXT.md) - Documentation maintenance guide

**Master Guide**: [`AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) - Read first for project overview

---

**Happy coding! 🚀 99% accuracy saves flocks.**

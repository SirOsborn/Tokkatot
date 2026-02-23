# Tokkatot 2.0: AI Service Specification

**Document Version**: 2.0  
**Last Updated**: February 2026  
**Status**: Production Ready  
**Framework**: PyTorch + FastAPI  
**Model**: Ensemble (EfficientNetB0 + DenseNet121)

---

## Overview

The Tokkatot AI Service provides real-time chicken disease detection through ensemble deep learning. It analyzes fecal images to identify diseases with 99% accuracy and safety-first voting mechanism to prioritize animal health.

**Key Features**:
- Ensemble voting system (dual models for reliability)
- 99% accuracy on test set (70,677 samples)
- REST API with FastAPI (modern, fast, auto-documentation)
- Real-time inference (2-5 seconds CPU)
- Docker containerized deployment
- Health checks & monitoring built-in

---

## Architecture

### Ensemble Model Design

```
Input Image (224x224)
    ↓
┌─────────────────────────────────────┐
│   EfficientNetB0 (98.05% recall)    │
└────────────┬────────────────────────┘
             │
             ├─→ Prediction: [0.98, 0.01, 0.01, 0.00]
             │
             ↓
        ┌──────────┐
        │ Ensemble │  92% Confidence
        │ Voting   │  "Healthy"
        └──────────┘
             ↑
             │
┌────────────┴────────────────────────┐
│   DenseNet121 (96.69% recall)       │
└─────────────────────────────────────┘
             │
             └─→ Prediction: [0.96, 0.02, 0.01, 0.01]
```

### Safety-First Logic

```
1. Get predictions from both models
2. If AGREEMENT (both models agree):
   - Use ensemble confidence
   - Normal action (isolate if disease)
3. If DISAGREEMENT:
   - Mark as UNCERTAIN
   - Flag for manual review
   - Conservative: ISOLATE for safety
```

### System Architecture

```
┌─────────────────────────────────────┐
│   Mobile / Web Client               │
│   (Take photo or upload)            │
└────────────────┬────────────────────┘
                 │ HTTPS
                 ↓
┌─────────────────────────────────────────────┐
│   FastAPI Service (Port 8000)               │
├─────────────────────────────────────────────┤
│  POST /predict                              │
│  GET /health                                │
│  POST /predict/detailed                     │
│  GET /docs (Swagger UI)                     │
└────────────────┬────────────────────────────┘
                 │
                 ↓
        ┌────────────────────┐
        │ Image Preprocessing│
        │ • Resize to 224x224│
        │ • Normalize RGB    │
        └────────────┬───────┘
                     │
         ┌───────────┴───────────┐
         ↓                       ↓
    ┌─────────────┐      ┌──────────────┐
    │ EfficientNet│      │ DenseNet121  │
    │ Checkpoint  │      │ Checkpoint   │
    └────────┬────┘      └────────┬─────┘
             │                    │
             ├─→ Softmax (4 classes)
             │   └─→ [Coccidiosis, Healthy, Newcastle, Salmonella]
             │
             ├──────────┤
             │ Voting   │
             ├──────────┤
             │
             ↓
    ┌───────────────────┐
    │ Safety Thresholds │
    │ • Healthy: 80%    │
    │ • Uncertain: 50%  │
    └────────┬──────────┘
             │
             ↓
    ┌─────────────────────┐
    │ Risk Assessment     │
    │ • Classification    │
    │ • Risk Level        │
    │ • Action Required   │
    └─────────────────────┘
```

---

## Supported Disease Classes

| Class | Detection | Treatment | Isolation |
|-------|-----------|-----------|-----------|
| **Healthy** | Normal droppings | None | No |
| **Coccidiosis** | Bloody/watery droppings | Anticoccidial drugs | Yes |
| **Newcastle Disease** | Green droppings, twisting neck | Supportive care, vaccines | Yes |
| **Salmonella** | Pale, whitish droppings | Antibiotics | Yes |

---

## API Specification

### Base URL
```
Production: https://api.tokkatot.local/ai-service
Development: http://localhost:8000
```

### Endpoints

#### 1. Health Check
```
GET /health

Response (200):
{
  "status": "healthy",
  "model_loaded": true,
  "device": "cpu",
  "model_version": "1.0.0",
  "ensemble_status": "ready"
}
```

#### 2. Disease Prediction
```
POST /predict
Content-Type: multipart/form-data

Request:
{
  "file": <binary image file>
}

Response (200):
{
  "classification": "Coccidiosis",
  "risk_level": "critical",
  "should_isolate": true,
  "action": "Isolate immediately. Contact veterinarian.",
  "confidence": 0.98,
  "model_votes": {
    "efficientnet": "Coccidiosis (0.99)",
    "densenet": "Coccidiosis (0.97)",
    "agreement": "unanimous"
  }
}
```

#### 3. Detailed Prediction Analysis
```
POST /predict/detailed
Content-Type: multipart/form-data

Request:
{
  "file": <binary image file>
}

Response (200):
{
  "classification": "Coccidiosis",
  "risk_level": "critical",
  "should_isolate": true,
  "action": "Isolate immediately. Contact veterinarian.",
  "confidence": 0.98,
  "class_probabilities": {
    "Healthy": 0.01,
    "Coccidiosis": 0.98,
    "Newcastle": 0.0,
    "Salmonella": 0.01
  },
  "model_analysis": {
    "efficientnet": {
      "prediction": "Coccidiosis",
      "confidence": 0.99,
      "probabilities": [0.01, 0.99, 0.0, 0.0]
    },
    "densenet": {
      "prediction": "Coccidiosis",
      "confidence": 0.97,
      "probabilities": [0.02, 0.97, 0.0, 0.01]
    }
  },
  "ensemble_decision": {
    "voting_method": "weighted_average",
    "agreement_level": "unanimous",
    "confidence_after_voting": 0.98
  },
  "processing_time_ms": 3200,
  "image_shape": [224, 224, 3]
}
```

---

## Performance Specifications

### Model Performance

| Metric | Score |
|--------|-------|
| **Overall Accuracy** | 99% |
| **Overall Recall** | 99% |
| **Overall Precision** | 99% |
| **F1 Score** | 99% |

### Component Models

| Model | Validation Recall | Architecture | Weights |
|-------|-------------------|--------------|---------|
| EfficientNetB0 | 98.05% | Lightweight CNN | 52 MB |
| DenseNet121 | 96.69% | Dense CNN | 33 MB |
| Ensemble | 99% | Voting mechanism | 47.2 MB (combined) |

### Inference Performance

| Metric | Performance |
|--------|-------------|
| **Prediction Time (CPU)** | 2-5 seconds |
| **Prediction Time (GPU)** | 0.5-1 second |
| **Memory Usage** | 2-3 GB (model + runtime) |
| **Model Size** | 47.2 MB |
| **Maximum Concurrent Requests** | 4-8 (depends on hardware) |

### Hardware Requirements

**Minimum (CPU only)**:
- CPU: 2 cores @ 1.5GHz+
- RAM: 2GB
- Storage: 100MB (model + runtime)

**Recommended (Best Performance)**:
- CPU: 4 cores @ 2.0GHz+
- RAM: 4GB
- GPU: NVIDIA GPU with CUDA (optional, significantly faster)

---

## Deployment

### Docker Container

**Image**: `tokkatot/ai-service:latest`

**Build**:
```bash
docker build -t tokkatot-ai:latest .
```

**Run**:
```bash
docker run -d \
  --name tokkatot-ai \
  -p 8000:8000 \
  -e PYTHONUNBUFFERED=1 \
  --restart unless-stopped \
  tokkatot/ai-service:latest
```

**Docker Compose**:
```yaml
services:
  tokkatot-ai:
    build:
      context: .
      dockerfile: Dockerfile
    image: tokkatot-ai:latest
    ports:
      - "8000:8000"
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "python", "-c", "import requests; requests.get('http://localhost:8000/health')"]
      interval: 30s
      timeout: 10s
      retries: 3
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G
```

### Model Files Location

Inside container:
```
/app/
├── app.py
├── inference.py
├── models.py
├── data_utils.py
└── outputs/
    ├── ensemble_model.pth (47.2 MB)
    └── checkpoints/
        ├── EfficientNetB0_best.pth
        └── DenseNet121_best.pth
```

---

## Integration with Tokkatot Backend

### Request Flow

```
1. Farmer takes photo with mobile device
   ↓
2. Mobile app sends to Go middleware (/api/ai/predict)
   ↓
3. Go middleware proxies to AI service (http://tokkatot-ai:8000/predict)
   ↓
4. AI service processes image
   ↓
5. Returns prediction with confidence
   ↓
6. Go middleware forwards to mobile app
   ↓
7. Mobile displays result with safety recommendations
```

### Integration Points

**Endpoint**: `POST /api/ai/predict` (Go Middleware)
- Receives image from mobile client
- Validates image (max 5MB, supported formats)
- Calls AI service: `POST http://tokkatot-ai:8000/predict`
- Returns result to mobile app

**Real-Time Processing**:
- Optional: WebSocket stream from mobile camera
- Allows continuous analysis with each frame
- Better for live monitoring during illness outbreak

---

## Error Handling

### Error Codes

| Code | Message | Cause | Action |
|------|---------|-------|--------|
| 400 | Invalid image format | File not image/JPEG/PNG | Retry with valid image |
| 413 | File too large | Image > 5MB | Reduce image size |
| 500 | Model not loaded | Startup failure | Restart service |
| 503 | Service unavailable | Processing queue full | Retry after 5 seconds |
| 504 | Request timeout | Inference > 30 seconds | Check hardware resources |

### Response Errors

```json
{
  "error": {
    "code": "INVALID_IMAGE_FORMAT",
    "message": "Uploaded file is not a valid image (JPEG or PNG required)",
    "details": "File must be JPEG, PNG, or WebP format"
  }
}
```

---

## Security

### Authentication
- AI service runs on internal network (not exposed to internet)
- Go middleware acts as security gateway
- JWT token validation in middleware layer

### Data Privacy
- Images not stored on AI service
- No logging of image content
- Processing in-memory only
- Results returned immediately

### Model Security
- Model weights in binary format (not interpretable)
- No model inversion attacks possible
- Read-only container filesystem

---

## Monitoring & Logging

### Health Checks
- Executed every 30 seconds
- Checks model availability
- Tracks inference latency
- Alerts on degradation

### Metrics to Monitor
- Request rate (pred/minute)
- Average inference time
- Error rate (failed predictions)
- Memory usage
- GPU utilization (if available)
- Queue depth

### Logging
```
[2026-02-19 15:30:45] INFO: Disease detected - Coccidiosis (confidence: 0.98)
[2026-02-19 15:30:50] INFO: Health check passed - Model loaded
[2026-02-19 15:31:05] ERROR: Failed to load image - Invalid format
```

---

## Testing

### Unit Tests
```bash
# Test individual models
python -m pytest tests/test_models.py

# Test inference pipeline
python -m pytest tests/test_inference.py

# Test API endpoints
python -m pytest tests/test_api.py
```

### Integration Tests
```bash
# Test with Go middleware
curl -X POST http://localhost:8000/predict \
  -F "file=@test_image.jpg"
```

### Performance Tests
```bash
# Measure inference time
python -m pytest tests/test_performance.py --benchmark

# Stress test (concurrent requests)
locust -f tests/stress_test.py
```

---

## Constraints & Limitations

### Image Constraints
- Format: JPEG, PNG, WebP
- Size: Maximum 5MB
- Aspect ratio: Any (will be resized to 224x224)
- Minimum resolution: 224x224 recommended

### Processing Constraints
- Batch size: 1 (single image per request)
- Maximum concurrent requests: 4-8
- Timeout: 30 seconds per request
- Queue depth: 100 requests max

### Accuracy Constraints
- Trained on chicken fecal images only
- NOT applicable to other animals
- Assumes good lighting conditions
- Performance degrads with poor image quality

---

## Future Enhancements

### Planned Improvements
1. **Batch Processing** - Process multiple images in one request
2. **Streaming Input** - Accept video stream from camera
3. **Model Retraining** - Periodic model updates with new data
4. **Uncertainty Quantification** - Detailed confidence metrics
5. **Export Functionality** - Allow model export for edge deployment
6. **Multi-language Output** - Disease names in Khmer/English

### Technical Debt
1. Optimize model size (15-20% reduction possible)
2. Add model versioning (track model updates)
3. Implement model serving with KServe
4. Add distributed inference (multiple GPU support)

---

## Related Documents

- [IG_SPECIFICATIONS_API.md](IG_SPECIFICATIONS_API.md) - API endpoints (includes AI endpoints)
- [01_SPECIFICATIONS_ARCHITECTURE.md](01_SPECIFICATIONS_ARCHITECTURE.md) - System architecture
- [OG_SPECIFICATIONS_DEPLOYMENT.md](OG_SPECIFICATIONS_DEPLOYMENT.md) - Deployment procedures
- [ai-service/README.md](../ai-service/README.md) - Quick start guide
- [ai-service/MODEL_CARD_ENSEMBLE.md](../ai-service/MODEL_CARD_ENSEMBLE.md) - Detailed model documentation

---

**Version History**

| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | PyTorch ensemble, FastAPI, production-ready |
| 1.0 | Jan 2026 | TensorFlow EfficientNetB0, Flask (prototype) |

**Status**: Ready for Integration with Tokkatot 2.0 Backend

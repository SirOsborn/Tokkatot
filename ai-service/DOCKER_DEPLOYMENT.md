# Tokkatot AI - Docker Deployment Guide

## Quick Start

### Using Docker Compose (Recommended)

1. **Build and run the container:**
   ```bash
   docker-compose up --build
   ```

2. **Access the API:**
   - API Documentation: http://localhost:8000/docs
   - Health Check: http://localhost:8000/health
   - Root: http://localhost:8000/

3. **Stop the container:**
   ```bash
   docker-compose down
   ```

### Using Docker Directly

1. **Build the image:**
   ```bash
   docker build -t tokkatot-ai:latest .
   ```

2. **Run the container:**
   ```bash
   docker run -d -p 8000:8000 --name tokkatot-ai tokkatot-ai:latest
   ```

3. **View logs:**
   ```bash
   docker logs -f tokkatot-ai
   ```

4. **Stop and remove:**
   ```bash
   docker stop tokkatot-ai
   docker rm tokkatot-ai
   ```

## API Endpoints

### 1. Health Check
```bash
curl http://localhost:8000/health
```

### 2. Simple Prediction
```bash
curl -X POST "http://localhost:8000/predict" \
  -H "Content-Type: multipart/form-data" \
  -F "file=@path/to/image.jpg"
```

Response:
```json
{
  "classification": "Healthy",
  "risk_level": "LOW RISK",
  "should_isolate": false,
  "action": "Clear for flock",
  "confidence": 0.95
}
```

### 3. Detailed Prediction
```bash
curl -X POST "http://localhost:8000/predict/detailed" \
  -H "Content-Type: multipart/form-data" \
  -F "file=@path/to/image.jpg"
```

Response includes individual model predictions, probabilities, and ensemble details.

### 4. Safety Evaluation
```bash
curl -X POST "http://localhost:8000/evaluate/safety" \
  -H "Content-Type: multipart/form-data" \
  -F "file=@path/to/image.jpg"
```

## Testing the API

### Using Python
```python
import requests

# Health check
response = requests.get("http://localhost:8000/health")
print(response.json())

# Predict
with open("test_image.jpg", "rb") as f:
    files = {"file": f}
    response = requests.post("http://localhost:8000/predict", files=files)
    print(response.json())
```

### Using PowerShell
```powershell
# Health check
Invoke-RestMethod -Uri "http://localhost:8000/health"

# Predict
$file = Get-Item "test_image.jpg"
$multipartContent = [System.Net.Http.MultipartFormDataContent]::new()
$fileStream = [System.IO.FileStream]::new($file.FullName, [System.IO.FileMode]::Open)
$fileContent = [System.Net.Http.StreamContent]::new($fileStream)
$multipartContent.Add($fileContent, "file", $file.Name)

$response = Invoke-RestMethod -Uri "http://localhost:8000/predict" -Method Post -Body $multipartContent -ContentType "multipart/form-data"
$response
```

## Cloud Deployment

### AWS ECS/Fargate
1. Push image to ECR:
   ```bash
   aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
   docker tag tokkatot-ai:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/tokkatot-ai:latest
   docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/tokkatot-ai:latest
   ```

2. Create ECS task definition and service using the pushed image.

### Google Cloud Run
```bash
# Build and push
gcloud builds submit --tag gcr.io/PROJECT-ID/tokkatot-ai

# Deploy
gcloud run deploy tokkatot-ai \
  --image gcr.io/PROJECT-ID/tokkatot-ai \
  --platform managed \
  --region us-central1 \
  --memory 4Gi \
  --cpu 2
```

### Azure Container Instances
```bash
# Push to ACR
az acr build --registry myregistry --image tokkatot-ai:latest .

# Deploy
az container create \
  --resource-group myResourceGroup \
  --name tokkatot-ai \
  --image myregistry.azurecr.io/tokkatot-ai:latest \
  --cpu 2 \
  --memory 4 \
  --port 8000 \
  --dns-name-label tokkatot-ai
```

### Docker Hub
```bash
# Tag and push
docker tag tokkatot-ai:latest yourusername/tokkatot-ai:latest
docker push yourusername/tokkatot-ai:latest
```

## Project Structure

The project uses `pyproject.toml` for dependency management. The Dockerfile automatically installs all dependencies defined in the project configuration.

## Configuration

### Environment Variables (optional)
You can customize the service by setting environment variables:

```yaml
environment:
  - MODEL_PATH=outputs/ensemble_model.pth
  - HEALTHY_THRESHOLD=0.80
  - UNCERTAINTY_THRESHOLD=0.50
  - PORT=8000
```

### Resource Requirements
- **Minimum:** 2GB RAM, 1 CPU
- **Recommended:** 4GB RAM, 2 CPUs
- **GPU:** Not required (optimized for CPU inference)

## Troubleshooting

### Container won't start
```bash
# Check logs
docker logs tokkatot-ai

# Check if port is already in use
netstat -ano | findstr :8000  # Windows
lsof -i :8000  # Linux/Mac
```

### Model not loading
- Ensure `outputs/ensemble_model.pth` exists
- Check model file permissions
- Verify sufficient memory (4GB recommended)

### Slow predictions
- Consider increasing CPU/memory allocation
- Model is optimized for CPU, but GPU can be enabled by modifying Dockerfile

## Development Mode

For development with hot reload:
```bash
docker run -d \
  -p 8000:8000 \
  -v ${PWD}:/app \
  --name tokkatot-ai-dev \
  tokkatot-ai:latest \
  uvicorn app:app --host 0.0.0.0 --port 8000 --reload
```

## Security Notes

- Container runs as non-root user (appuser)
- Only necessary ports exposed
- For production: Configure CORS in [app.py](app.py) appropriately
- For production: Add authentication/API keys
- For production: Enable HTTPS/TLS

## Monitoring

The container includes a health check endpoint that monitors model status:
```bash
curl http://localhost:8000/health
```

Response:
```json
{
  "status": "healthy",
  "model_loaded": true,
  "device": "cpu"
}
```

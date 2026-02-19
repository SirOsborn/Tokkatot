# EfficientNetB0 - Chicken Disease Detection Model

**Model Version:** v1.0.0  
**Release Date:** January 16, 2026  
**© 2026 Tokkatot. All Rights Reserved.**

## 📊 Model Performance

**Training completed:** 90 epochs  
**Dataset:** 400,000 training images, 40,000 validation images

### Metrics (Validation Set)

| Metric | Score |
|--------|-------|
| **Overall Recall** | **98.05%** |
| **Overall Accuracy** | **98.05%** |
| **F1 Score** | **98.05%** |

### Per-Class Recall

| Disease Class | Recall |
|---------------|--------|
| Coccidiosis | 98.98% |
| Healthy | 96.58% |
| New Castle Disease | 98.18% |
| Salmonella | 98.47% |

## 🎯 Model Details

- **Architecture:** EfficientNetB0 (pretrained on ImageNet, fine-tuned)
- **Input Size:** 224x224 RGB images
- **Output Classes:** 4 (Coccidiosis, Healthy, New Castle Disease, Salmonella)
- **Framework:** PyTorch 2.5.1
- **Optimization:** AdamW with ReduceLROnPlateau
- **Loss Function:** RecallFocusedLoss (5x false negative penalty)

## 🚀 Usage

### Installation

```bash
# Install PyTorch with CUDA support
pip install torch torchvision --index-url https://download.pytorch.org/whl/cu121

# Install other dependencies
pip install pillow numpy
```

### Loading the Model

```python
import torch
from models import EfficientNetB0Classifier

# Initialize model
model = EfficientNetB0Classifier(num_classes=4)

# Load checkpoint
checkpoint = torch.load('EfficientNetB0_best.pth', map_location='cpu')
model.load_state_dict(checkpoint['model_state_dict'])
model.eval()

# Move to GPU if available
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
model = model.to(device)
```

### Inference

```python
from PIL import Image
import torchvision.transforms as transforms

# Define transforms (same as training)
transform = transforms.Compose([
    transforms.Resize((224, 224)),
    transforms.ToTensor(),
    transforms.Normalize(
        mean=[0.485, 0.456, 0.406],
        std=[0.229, 0.224, 0.225]
    )
])

# Load and preprocess image
image = Image.open('fecal_image.jpg').convert('RGB')
input_tensor = transform(image).unsqueeze(0).to(device)

# Predict
with torch.no_grad():
    outputs = model(input_tensor)
    probabilities = torch.nn.functional.softmax(outputs, dim=1)
    confidence, predicted = torch.max(probabilities, 1)

# Class mapping
CLASS_NAMES = ['Coccidiosis', 'Healthy', 'New Castle Disease', 'Salmonella']
predicted_class = CLASS_NAMES[predicted.item()]
confidence_score = confidence.item()

print(f"Prediction: {predicted_class}")
print(f"Confidence: {confidence_score:.2%}")
```

## 📦 Files Needed

To use this model, you need:
1. **EfficientNetB0_best.pth** - Model checkpoint (this file)
2. **models.py** - Model architecture definition
3. **PIL, torch, torchvision** - Python dependencies

## ⚠️ Important Notes

### Safety Recommendations

This model achieves 98.05% recall, but for production deployment:
- **Use ensemble voting** with DenseNet121 model for higher reliability
- **Set confidence thresholds** (recommended: 80% for healthy classification)
- **Implement isolation protocol** for uncertain predictions
- **Never use as sole diagnostic tool** - always validate with veterinary assessment

### Limitations

- Trained on specific dataset of chicken fecal images
- Performance may vary with different lighting, camera angles, or image quality
- Not validated for other poultry species
- Requires proper image preprocessing (resize, normalize)

## 🔧 Hardware Requirements

### Minimum
- **CPU:** Any modern processor (inference will be slower)
- **RAM:** 2GB free memory
- **Storage:** 100MB

### Recommended
- **GPU:** NVIDIA GPU with CUDA support (RTX 3060 or better)
- **VRAM:** 2GB+ for batch inference
- **RAM:** 4GB+ free memory

### Edge Deployment (Raspberry Pi)
- Compatible with Raspberry Pi 4 (4GB+ RAM)
- Inference time: ~2-5 seconds per image on CPU
- Consider model quantization for faster inference

## 📄 License

**© 2026 Tokkatot. All Rights Reserved.**

This model is proprietary and part of Tokkatot Smart Chicken Farming Solutions. Commercial use requires explicit permission from Tokkatot.

## 📧 Contact

**Tokkatot Smart Chicken Farming Solutions**
- Email: tokkatot.info@gmail.com
- Website: tokkatot.aztrolabe.com
- AI Engineer: sunhenglong@outlook.com

## 📝 Changelog

### v1.0.0 (2026-01-16)
- Initial release
- 90 epochs training completed
- 98.05% validation recall achieved
- 4-class classification (Coccidiosis, Healthy, New Castle Disease, Salmonella)

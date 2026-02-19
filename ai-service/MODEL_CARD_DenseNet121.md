# DenseNet121 - Chicken Disease Detection Model

**Model Version:** v1.0.0  
**Release Date:** January 17, 2026  
**© 2026 Tokkatot. All Rights Reserved.**

## 📊 Model Performance

**Training completed:** 20 epochs  
**Dataset:** 400,000 training images, 40,000 validation images

### Metrics (Validation Set)

| Metric | Score |
|--------|-------|
| **Overall Recall** | **96.69%** |
| **Overall Accuracy** | **96.69%** |
| **F1 Score** | **96.69%** |

### Per-Class Recall

| Disease Class | Recall |
|---------------|--------|
| Coccidiosis | 98.08% |
| Healthy | 93.39% |
| New Castle Disease | 97.84% |
| Salmonella | 97.46% |

### Test Set Performance (Ensemble)

| Disease Class | Precision | Recall | F1-Score | Support |
|---------------|-----------|--------|----------|---------|
| Coccidiosis | 0.99 | 1.00 | 0.99 | 18,338 |
| Healthy | 1.00 | 0.98 | 0.99 | 15,451 |
| New Castle Disease | 0.99 | 1.00 | 0.99 | 15,339 |
| Salmonella | 0.99 | 1.00 | 1.00 | 18,009 |

**Test Accuracy:** 99% (67,137 classified samples)

## 🎯 Model Details

- **Architecture:** DenseNet121 (pretrained on ImageNet, fine-tuned)
- **Input Size:** 224x224 RGB images
- **Output Classes:** 4 (Coccidiosis, Healthy, New Castle Disease, Salmonella)
- **Framework:** PyTorch 2.5.1
- **Optimization:** AdamW with ReduceLROnPlateau
- **Loss Function:** RecallFocusedLoss (5x false negative penalty)
- **Key Features:**
  - Dense connections for superior feature reuse
  - Robust gradient flow for stable training
  - Excellent at capturing fine-grained patterns
  - Complementary to EfficientNetB0 in ensemble

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
from models import DenseNet121Classifier

# Initialize model
model = DenseNet121Classifier(num_classes=4)

# Load checkpoint
checkpoint = torch.load('DenseNet121_best.pth', map_location='cpu')
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
1. **DenseNet121_best.pth** - Model checkpoint (this file)
2. **models.py** - Model architecture definition
3. **PIL, torch, torchvision** - Python dependencies

## ⚠️ Important Notes

### Safety Recommendations

This model achieves 96.69% recall, but for production deployment:
- **Use ensemble voting** with EfficientNetB0 model for higher reliability
- **Set confidence thresholds** (recommended: 80% for healthy classification)
- **Implement isolation protocol** for uncertain predictions
- **Never use as sole diagnostic tool** - always validate with veterinary assessment

### Ensemble Benefits

DenseNet121 complements EfficientNetB0 by:
- Providing a second opinion for safety-critical decisions
- Capturing different feature patterns through dense connectivity
- Enabling uncertainty detection through model disagreement
- Achieving **99% test accuracy** when combined in ensemble

### Limitations

- Trained on specific dataset of chicken fecal images
- Performance may vary with different lighting, camera angles, or image quality
- Not validated for other poultry species
- Requires proper image preprocessing (resize, normalize)
- Computationally more intensive than EfficientNetB0

## 🔧 Hardware Requirements

### Minimum
- **CPU:** Any modern processor (inference will be slower)
- **RAM:** 3GB free memory
- **Storage:** 150MB

### Recommended
- **GPU:** NVIDIA GPU with CUDA support (RTX 3060 or better)
- **VRAM:** 4GB+ for batch inference
- **RAM:** 8GB+ free memory

### Deployment Notes
- DenseNet121 is more compute-intensive than EfficientNetB0
- For edge devices (Raspberry Pi), use EfficientNetB0 alone
- For server/cloud deployment, use full ensemble for maximum accuracy

## 🏆 Performance Highlights

### Training Progress

- **Epoch 1:** 87.16% recall → **Epoch 20:** 96.69% recall
- Steady improvement with early stopping at epoch 20
- Best model saved based on validation recall
- No overfitting observed (train/val metrics aligned)

### Why DenseNet121?

1. **Dense Connectivity:** Each layer receives features from all preceding layers
2. **Parameter Efficiency:** Fewer parameters than ResNet while maintaining accuracy
3. **Strong Gradients:** Alleviates vanishing gradient problem
4. **Feature Reuse:** Excellent for medical/diagnostic imaging tasks

## 📄 License

**© 2026 Tokkatot. All Rights Reserved.**

This model is proprietary and part of Tokkatot Smart Chicken Farming Solutions. Commercial use requires explicit permission from Tokkatot.

## 📧 Contact

**Tokkatot Smart Chicken Farming Solutions**
- Email: tokkatot.info@gmail.com
- Website: tokkatot.aztrolabe.com
- AI Engineer: sunhenglong@outlook.com

## 📝 Changelog

### v1.0.0 (2026-01-17)
- Initial release
- 20 epochs training completed
- 96.69% validation recall achieved
- 4-class classification (Coccidiosis, Healthy, New Castle Disease, Salmonella)
- Integrated into ensemble system with EfficientNetB0

---

**Part of the Tokkatot Safety-First Ensemble AI System**

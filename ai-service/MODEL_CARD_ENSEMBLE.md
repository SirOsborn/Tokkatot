# Ensemble Model - Chicken Disease Detection System

**Model Version:** v1.0.0  
**Release Date:** January 17, 2026  
**© 2026 Tokkatot. All Rights Reserved.**

## 🎯 Overview

The Tokkatot Ensemble Model combines two state-of-the-art deep learning architectures into a single safety-first system for chicken disease detection. By leveraging the complementary strengths of **EfficientNetB0** and **DenseNet121**, the ensemble achieves **99% accuracy** on real-world test data while maintaining maximum safety through a parallel voting mechanism.

## 📊 Performance Summary

### Test Set Performance

**Total Samples:** 70,677  
**Classified:** 67,137 (94.99%)  
**Isolated for Safety:** 3,540 (5.01%)

| Metric | Score |
|--------|-------|
| **Overall Accuracy** | **99%** |
| **Overall Recall** | **99%** |
| **Overall Precision** | **99%** |
| **F1 Score** | **99%** |

### Per-Class Performance

| Disease Class | Precision | Recall | F1-Score | Support |
|---------------|-----------|--------|----------|---------|
| **Coccidiosis** | 0.99 | 1.00 | 0.99 | 18,338 |
| **Healthy** | 1.00 | 0.98 | 0.99 | 15,451 |
| **New Castle Disease** | 0.99 | 1.00 | 0.99 | 15,339 |
| **Salmonella** | 0.99 | 1.00 | 1.00 | 18,009 |

### Component Model Performance

| Model | Validation Recall | Epochs Trained |
|-------|-------------------|----------------|
| **EfficientNetB0** | 98.05% | 90 |
| **DenseNet121** | 96.69% | 20 |
| **Ensemble** | **99%** | - |

## 🏗️ Architecture

### Dual-Model Parallel Voting System

```
                    ┌───────────────────────┐
                    │    Input Image        │
                    │   (224x224 RGB)       │
                    └───────────┬───────────┘
                                │
                    ┌───────────┴───────────┐
                    │                       │
         ┌──────────▼──────────┐  ┌────────▼─────────┐
         │   EfficientNetB0    │  │   DenseNet121    │
         │   (98.05% recall)   │  │  (96.69% recall) │
         └──────────┬──────────┘  └────────┬─────────┘
                    │                       │
                    │   ┌───────────────┐   │
                    └───► Safety Vote   ◄───┘
                        │    Logic      │
                        └───────┬───────┘
                                │
                    ┌───────────▼───────────┐
                    │    Classification     │
                    │  or Isolation Action  │
                    └───────────────────────┘
```

### Safety-First Decision Logic

The system **isolates** a chicken if **ANY** of the following conditions are met:

1. **Uncertainty Check**
   - Either model's maximum confidence < 50%
   - Action: `ISOLATE` (unknown/out-of-distribution sample)

2. **Safety Vote (Healthy Classification)**
   - Either model's healthy confidence < 80%
   - Action: `ISOLATE` (potential disease detected)

3. **Model Disagreement**
   - Models predict different classes AND either predicts disease
   - Action: `ISOLATE` (safety override on disagreement)

4. **Classification (Safe Cases)**
   - Both models confident and agree → Return predicted class
   - Both predict healthy with high confidence → `Healthy`

## 🎯 Key Features

### Safety Mechanisms

- ✅ **Parallel Voting:** Two independent models provide redundancy
- ✅ **Conservative Thresholds:** 80% confidence required for healthy classification
- ✅ **Uncertainty Detection:** Isolate low-confidence predictions
- ✅ **Disagreement Handling:** Isolate when models disagree on disease cases
- ✅ **No False Negatives:** Prioritizes catching all diseased birds

### Technical Highlights

- **Multi-Model Ensemble:** Combines CNN architectures with different inductive biases
- **Complementary Strengths:**
  - EfficientNetB0: Fast, efficient, optimized for general features
  - DenseNet121: Dense connections, superior feature reuse, fine-grained patterns
- **Production Ready:** Single `.pth` file contains both models
- **Configurable Thresholds:** Adjust safety levels based on farm requirements

## 🚀 Usage

### Installation

```bash
# Install dependencies
pip install torch torchvision pillow numpy

# Or use the project setup
cd tokkatot_ai
pip install -e .
```

### Quick Start

```python
from inference import ChickenDiseaseDetector

# Initialize ensemble detector
detector = ChickenDiseaseDetector(
    model_path='outputs/ensemble_model.pth',
    healthy_threshold=0.80,      # 80% confidence for healthy
    uncertainty_threshold=0.50    # 50% minimum confidence
)

# Single prediction
result = detector.predict('fecal_image.jpg')
print(f"Classification: {result}")
# Output: 'Healthy', 'Salmonella', 'Coccidiosis', 'New Castle Disease', or 'ISOLATE'
```

### Detailed Predictions

```python
# Get full prediction details
result = detector.predict('fecal_image.jpg', return_details=True)

print(f"Classification: {result['classification']}")
print(f"Should Isolate: {result['should_isolate']}")
print(f"Isolation Reason: {result['isolation_reason']}")
print(f"Action: {result['action']}")

# Individual model predictions
print("\nEfficientNetB0:")
print(f"  Prediction: {result['models']['efficientnet']['prediction']}")
print(f"  Confidence: {result['models']['efficientnet']['confidence']:.2%}")
print(f"  Probabilities: {result['models']['efficientnet']['probabilities']}")

print("\nDenseNet121:")
print(f"  Prediction: {result['models']['densenet']['prediction']}")
print(f"  Confidence: {result['models']['densenet']['confidence']:.2%}")
print(f"  Probabilities: {result['models']['densenet']['probabilities']}")
```

### Batch Processing

```python
# Process multiple images
image_paths = ['img1.jpg', 'img2.jpg', 'img3.jpg']
results = detector.predict_batch(image_paths, return_details=True)

for img_path, result in zip(image_paths, results):
    print(f"{img_path}: {result['classification']}")
```

### Safety Evaluation

```python
# Check if an image requires isolation
is_safe, reason = detector.evaluate_safety('fecal_image.jpg')

if not is_safe:
    print(f"⚠️ ISOLATE CHICKEN - Reason: {reason}")
else:
    print(f"✓ Safe to proceed - {reason}")
```

## 🔧 Configuration

### Adjusting Safety Thresholds

```python
# More conservative (more isolations)
detector = ChickenDiseaseDetector(
    model_path='outputs/ensemble_model.pth',
    healthy_threshold=0.90,      # Require 90% confidence for healthy
    uncertainty_threshold=0.60    # Require 60% minimum confidence
)

# Less conservative (fewer isolations)
detector = ChickenDiseaseDetector(
    model_path='outputs/ensemble_model.pth',
    healthy_threshold=0.70,      # Require 70% confidence for healthy
    uncertainty_threshold=0.40    # Require 40% minimum confidence
)
```

### Threshold Guidelines

| Threshold | Value Range | Impact |
|-----------|-------------|--------|
| `healthy_threshold` | 0.70-0.90 | Higher = More isolations for safety |
| `uncertainty_threshold` | 0.40-0.70 | Higher = Reject more uncertain predictions |

**Recommended Settings:**
- **High-risk farm:** `healthy_threshold=0.90`, `uncertainty_threshold=0.60`
- **Standard farm:** `healthy_threshold=0.80`, `uncertainty_threshold=0.50` (default)
- **Low-risk farm:** `healthy_threshold=0.70`, `uncertainty_threshold=0.40`

## 📦 Model Files

### What's Included

The `ensemble_model.pth` file contains:
- ✅ EfficientNetB0 model weights (best checkpoint, 98.05% recall)
- ✅ DenseNet121 model weights (best checkpoint, 96.69% recall)
- ✅ Class mapping and metadata
- ✅ Configuration parameters

### Required Dependencies

```python
# Core dependencies
torch >= 2.0.0
torchvision >= 0.15.0
pillow >= 10.0.0
numpy >= 1.24.0

# Optional (for development)
scikit-learn >= 1.3.0
tqdm >= 4.65.0
tensorboard >= 2.13.0
```

## 🎓 Training Details

### Dataset

- **Training samples:** 400,000 images
- **Validation samples:** 40,000 images
- **Test samples:** 70,677 images
- **Classes:** 4 (Coccidiosis, Healthy, New Castle Disease, Salmonella)
- **Image size:** 224x224 pixels

### Training Configuration

- **Batch size:** 32
- **Optimizer:** AdamW (lr=1e-4)
- **Loss function:** RecallFocusedLoss (5x false negative penalty)
- **Scheduler:** ReduceLROnPlateau (patience=5)
- **Early stopping:** Patience=10 on validation recall
- **Data augmentation:** Rotation, flip, color jitter, affine transforms

### Why This Ensemble Works

1. **Complementary Architectures:**
   - EfficientNetB0: Compound scaling, efficient feature extraction
   - DenseNet121: Dense connections, gradient flow, feature reuse

2. **Different Training Dynamics:**
   - EfficientNetB0: Trained 90 epochs (early stopped)
   - DenseNet121: Trained 20 epochs (early stopped)
   - Different convergence patterns capture different aspects

3. **Recall-Focused Training:**
   - Both models optimized for maximum recall
   - False negative penalty ensures disease detection
   - Class weighting addresses imbalance

## 🏆 Performance Comparison

### Individual vs. Ensemble

| Model | Accuracy | Recall | Key Strength |
|-------|----------|--------|--------------|
| EfficientNetB0 | 98.05% | 98.05% | Speed, efficiency |
| DenseNet121 | 96.69% | 96.69% | Feature reuse, stability |
| **Ensemble** | **99%** | **99%** | **Safety, redundancy** |

### Real-World Benefits

- **5.01% isolation rate:** Manageable for farm operations
- **99% accuracy on classified samples:** Minimizes false positives
- **Near-perfect disease recall:** Catches virtually all diseased birds
- **Redundant system:** If one model fails, the other provides backup

## ⚠️ Important Notes

### Production Recommendations

1. **Use ensemble in production:** Don't rely on individual models alone
2. **Monitor isolation rate:** If >15%, consider adjusting thresholds
3. **Validate isolated birds:** Veterinary confirmation recommended
4. **Regular retraining:** Update models with new data periodically
5. **Hardware requirements:** GPU recommended for real-time inference

### Limitations

- Requires proper image quality and preprocessing
- Trained on specific chicken breeds and conditions
- Not validated for other poultry species
- Requires stable lighting and camera positioning
- Should not replace veterinary diagnosis

### Edge Deployment

For resource-constrained devices (Raspberry Pi):
- Use **EfficientNetB0 only** (faster inference)
- Full ensemble requires GPU for real-time performance
- Consider model quantization for edge optimization

## 🔧 Hardware Requirements

### Development/Training
- **GPU:** NVIDIA RTX 3060 or better (8GB+ VRAM)
- **RAM:** 16GB+ system memory
- **Storage:** 10GB+ free space

### Production/Inference
- **Server Deployment:** NVIDIA GPU (4GB+ VRAM) recommended
- **Edge Deployment:** Raspberry Pi 4 (4GB RAM) with EfficientNetB0 only
- **Cloud Deployment:** Standard GPU instances (AWS p3, GCP GPU)

### Inference Speed

| Hardware | Ensemble Inference Time |
|----------|-------------------------|
| RTX 3060 GPU | ~50ms per image |
| CPU (Intel i7) | ~2-3s per image |
| Raspberry Pi 4 | ~5-8s (EfficientNetB0 only) |

## 📄 License

**© 2026 Tokkatot. All Rights Reserved.**

This ensemble model is proprietary and part of Tokkatot Smart Chicken Farming Solutions. Commercial use, redistribution, or modification requires explicit written permission from Tokkatot.

## 📧 Contact

**Tokkatot Smart Chicken Farming Solutions**
- **Email:** tokkatot.info@gmail.com
- **Website:** tokkatot.aztrolabe.com
- **AI Engineer:** sunhenglong@outlook.com

## 📝 Changelog

### v1.0.0 (2026-01-17)
- Initial release of ensemble system
- Combined EfficientNetB0 (v1.0.0) + DenseNet121 (v1.0.0)
- Achieved 99% test accuracy
- 5.01% isolation rate (3,540/70,677 samples)
- Parallel safety vote mechanism implemented
- Configurable thresholds for farm-specific requirements

---

**Safety-First AI for Chicken Disease Detection**  
**Part of Tokkatot Smart Chicken Farming Solutions**

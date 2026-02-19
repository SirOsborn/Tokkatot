"""
Model architectures for chicken disease detection ensemble.
Combines EfficientNetB0 (fast, edge-friendly) and DenseNet (feature reuse, robust).
"""

import torch
import torch.nn as nn
from torchvision import models


class EfficientNetB0Classifier(nn.Module):
    """
    EfficientNetB0 for chicken disease classification.
    Optimized for edge deployment (Raspberry Pi) while maintaining accuracy.
    """
    
    def __init__(self, num_classes=4, dropout=0.3):
        super(EfficientNetB0Classifier, self).__init__()
        
        # Load pretrained EfficientNetB0
        self.backbone = models.efficientnet_b0(weights=models.EfficientNet_B0_Weights.IMAGENET1K_V1)
        
        # Modify classifier head for our 4 classes
        in_features = self.backbone.classifier[1].in_features
        self.backbone.classifier = nn.Sequential(
            nn.Dropout(p=dropout),
            nn.Linear(in_features, 512),
            nn.ReLU(),
            nn.Dropout(p=dropout/2),
            nn.Linear(512, num_classes)
        )
        
    def forward(self, x):
        return self.backbone(x)


class DenseNet121Classifier(nn.Module):
    """
    DenseNet121 for chicken disease classification.
    Excellent feature reuse and gradient flow for robust predictions.
    """
    
    def __init__(self, num_classes=4, dropout=0.3):
        super(DenseNet121Classifier, self).__init__()
        
        # Load pretrained DenseNet121
        self.backbone = models.densenet121(weights=models.DenseNet121_Weights.IMAGENET1K_V1)
        
        # Modify classifier head for our 4 classes
        in_features = self.backbone.classifier.in_features
        self.backbone.classifier = nn.Sequential(
            nn.Dropout(p=dropout),
            nn.Linear(in_features, 512),
            nn.ReLU(),
            nn.Dropout(p=dropout/2),
            nn.Linear(512, num_classes)
        )
        
    def forward(self, x):
        return self.backbone(x)


class EnsembleModel(nn.Module):
    """
    Ensemble wrapper for safety-first parallel inference.
    Both models run independently and vote on final decision.
    """
    
    def __init__(self, efficientnet_model, densenet_model):
        super(EnsembleModel, self).__init__()
        self.efficientnet = efficientnet_model
        self.densenet = densenet_model
        
    def forward(self, x):
        """
        Returns logits from both models separately.
        """
        efficient_out = self.efficientnet(x)
        dense_out = self.densenet(x)
        return efficient_out, dense_out
    
    def predict_with_safety_vote(self, x, healthy_threshold=0.80, uncertainty_threshold=0.50):
        """
        Safety-first prediction logic:
        1. Both models predict independently
        2. If EITHER model's healthy confidence < 80% → ISOLATE
        3. If ANY model's max confidence < 50% → ISOLATE (unknown/OOD)
        
        Args:
            x: Input tensor
            healthy_threshold: Confidence threshold for healthy class (default 0.80)
            uncertainty_threshold: Minimum confidence for any prediction (default 0.50)
        
        Returns:
            dict with:
                - prediction: Final prediction (0=Coccidiosis, 1=Healthy, 2=New Castle, 3=Salmonella, -1=ISOLATE)
                - should_isolate: Boolean flag
                - efficientnet_probs: Probability distribution from EfficientNet
                - densenet_probs: Probability distribution from DenseNet
                - confidence_scores: Dict with individual confidences
        """
        self.eval()
        with torch.no_grad():
            efficient_logits, dense_logits = self.forward(x)
            
            # Convert to probabilities
            efficient_probs = torch.softmax(efficient_logits, dim=1)
            dense_probs = torch.softmax(dense_logits, dim=1)
            
            # Get predictions and confidences
            efficient_conf, efficient_pred = torch.max(efficient_probs, dim=1)
            dense_conf, dense_pred = torch.max(dense_probs, dim=1)
            
            # Healthy is class index 1
            healthy_idx = 1
            efficient_healthy_conf = efficient_probs[:, healthy_idx]
            dense_healthy_conf = dense_probs[:, healthy_idx]
            
            # Safety checks
            should_isolate = torch.zeros(x.size(0), dtype=torch.bool, device=x.device)
            
            # Check 1: Uncertainty detection (max confidence < 50%)
            uncertainty_mask = (efficient_conf < uncertainty_threshold) | (dense_conf < uncertainty_threshold)
            should_isolate |= uncertainty_mask
            
            # Check 2: Safety vote on healthy predictions
            # If EITHER model is not confident about healthy (< 80%), isolate
            low_healthy_confidence = (efficient_healthy_conf < healthy_threshold) | (dense_healthy_conf < healthy_threshold)
            
            # For samples where both predict healthy but at least one isn't confident
            both_predict_healthy = (efficient_pred == healthy_idx) & (dense_pred == healthy_idx)
            should_isolate |= (both_predict_healthy & low_healthy_confidence)
            
            # Check 3: If models disagree and either predicts disease, isolate
            models_disagree = (efficient_pred != dense_pred)
            either_predicts_disease = (efficient_pred != healthy_idx) | (dense_pred != healthy_idx)
            should_isolate |= (models_disagree & either_predicts_disease)
            
            # Final prediction: use majority vote when both agree, else isolate
            final_pred = torch.where(
                should_isolate,
                torch.full_like(efficient_pred, -1),  # -1 = ISOLATE
                torch.where(
                    efficient_pred == dense_pred,
                    efficient_pred,  # Agreement
                    torch.full_like(efficient_pred, -1)  # Disagreement → ISOLATE
                )
            )
            
            return {
                'prediction': final_pred,
                'should_isolate': should_isolate,
                'efficientnet_probs': efficient_probs,
                'densenet_probs': dense_probs,
                'confidence_scores': {
                    'efficientnet_max': efficient_conf,
                    'densenet_max': dense_conf,
                    'efficientnet_healthy': efficient_healthy_conf,
                    'densenet_healthy': dense_healthy_conf
                }
            }


def create_ensemble(num_classes=4, pretrained=True):
    """
    Factory function to create the ensemble model.
    """
    efficientnet = EfficientNetB0Classifier(num_classes=num_classes)
    densenet = DenseNet121Classifier(num_classes=num_classes)
    ensemble = EnsembleModel(efficientnet, densenet)
    return ensemble

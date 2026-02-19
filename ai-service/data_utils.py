"""
Data utilities for loading and preprocessing chicken fecal images.
Handles the dataset structure with train/val/test splits.
"""

import os
from pathlib import Path
from typing import Tuple, Optional

import torch
from torch.utils.data import Dataset, DataLoader
from torchvision import transforms
from PIL import Image
import numpy as np


# Class mapping
CLASS_NAMES = ['Coccidiosis', 'Healthy', 'New Castle Disease', 'Salmonella']
CLASS_TO_IDX = {name: idx for idx, name in enumerate(CLASS_NAMES)}
IDX_TO_CLASS = {idx: name for name, idx in CLASS_TO_IDX.items()}


class ChickenFecalDataset(Dataset):
    """
    Dataset for chicken disease detection from fecal images.
    """
    
    def __init__(self, root_dir: str, split: str = 'train', transform=None):
        """
        Args:
            root_dir: Root directory containing train/val/test folders
            split: One of 'train', 'val', or 'test'
            transform: Optional transform to be applied on images
        """
        self.root_dir = Path(root_dir)
        self.split = split
        self.transform = transform
        
        # Build file list
        self.samples = []
        split_dir = self.root_dir / split
        
        if not split_dir.exists():
            raise ValueError(f"Split directory not found: {split_dir}")
        
        for class_name in CLASS_NAMES:
            class_dir = split_dir / class_name
            if not class_dir.exists():
                print(f"Warning: Class directory not found: {class_dir}")
                continue
            
            class_idx = CLASS_TO_IDX[class_name]
            
            # Get all image files
            for img_path in class_dir.glob('*'):
                if img_path.suffix.lower() in ['.jpg', '.jpeg', '.png', '.bmp']:
                    self.samples.append((str(img_path), class_idx))
        
        if len(self.samples) == 0:
            raise ValueError(f"No images found in {split_dir}")
        
        print(f"Loaded {len(self.samples)} images for {split} split")
        
        # Count samples per class
        class_counts = {}
        for _, class_idx in self.samples:
            class_name = IDX_TO_CLASS[class_idx]
            class_counts[class_name] = class_counts.get(class_name, 0) + 1
        
        print(f"Class distribution for {split}:")
        for class_name, count in class_counts.items():
            print(f"  {class_name}: {count}")
    
    def __len__(self):
        return len(self.samples)
    
    def __getitem__(self, idx):
        img_path, label = self.samples[idx]
        
        # Load image
        image = Image.open(img_path).convert('RGB')
        
        # Apply transforms
        if self.transform:
            image = self.transform(image)
        
        return image, label


def get_transforms(split: str = 'train', img_size: int = 224):
    """
    Get appropriate transforms for each split.
    Training includes augmentation for robustness.
    """
    if split == 'train':
        return transforms.Compose([
            transforms.Resize((img_size + 32, img_size + 32)),
            transforms.RandomCrop(img_size),
            transforms.RandomHorizontalFlip(p=0.5),
            transforms.RandomVerticalFlip(p=0.3),
            transforms.RandomRotation(20),
            transforms.ColorJitter(brightness=0.2, contrast=0.2, saturation=0.2, hue=0.1),
            transforms.RandomAffine(degrees=0, translate=(0.1, 0.1), scale=(0.9, 1.1)),
            transforms.ToTensor(),
            transforms.Normalize(mean=[0.485, 0.456, 0.406], 
                               std=[0.229, 0.224, 0.225])
        ])
    else:
        return transforms.Compose([
            transforms.Resize((img_size, img_size)),
            transforms.ToTensor(),
            transforms.Normalize(mean=[0.485, 0.456, 0.406], 
                               std=[0.229, 0.224, 0.225])
        ])


def create_dataloaders(
    data_dir: str,
    batch_size: int = 32,
    num_workers: int = 4,
    img_size: int = 224
) -> Tuple[DataLoader, DataLoader, DataLoader]:
    """
    Create train, validation, and test dataloaders.
    
    Args:
        data_dir: Root directory containing train/val/test folders
        batch_size: Batch size for training
        num_workers: Number of worker processes for data loading
        img_size: Size to resize images to
    
    Returns:
        Tuple of (train_loader, val_loader, test_loader)
    """
    # Create datasets
    train_dataset = ChickenFecalDataset(
        data_dir, 
        split='train',
        transform=get_transforms('train', img_size)
    )
    
    val_dataset = ChickenFecalDataset(
        data_dir,
        split='val',
        transform=get_transforms('val', img_size)
    )
    
    test_dataset = ChickenFecalDataset(
        data_dir,
        split='test',
        transform=get_transforms('test', img_size)
    )
    
    # Calculate class weights for handling imbalance (focus on disease detection)
    class_counts = np.zeros(len(CLASS_NAMES))
    for _, label in train_dataset.samples:
        class_counts[label] += 1
    
    # Inverse frequency weighting with extra emphasis on disease classes
    class_weights = 1.0 / (class_counts + 1e-6)
    
    # Boost disease classes even more (2x weight)
    # Healthy is index 1, so boost all others
    for idx in range(len(class_weights)):
        if idx != 1:  # Not healthy
            class_weights[idx] *= 2.0
    
    class_weights = class_weights / class_weights.sum() * len(class_weights)
    
    print(f"\nClass weights (disease detection focused):")
    for idx, weight in enumerate(class_weights):
        print(f"  {IDX_TO_CLASS[idx]}: {weight:.4f}")
    
    # Create dataloaders
    train_loader = DataLoader(
        train_dataset,
        batch_size=batch_size,
        shuffle=True,
        num_workers=num_workers,
        pin_memory=True,
        drop_last=True
    )
    
    val_loader = DataLoader(
        val_dataset,
        batch_size=batch_size,
        shuffle=False,
        num_workers=num_workers,
        pin_memory=True
    )
    
    test_loader = DataLoader(
        test_dataset,
        batch_size=batch_size,
        shuffle=False,
        num_workers=num_workers,
        pin_memory=True
    )
    
    return train_loader, val_loader, test_loader, torch.tensor(class_weights, dtype=torch.float32)


def calculate_class_weights(dataset: Dataset) -> torch.Tensor:
    """
    Calculate class weights for imbalanced dataset.
    Extra focus on not missing disease cases.
    """
    class_counts = np.zeros(len(CLASS_NAMES))
    
    for _, label in dataset.samples:
        class_counts[label] += 1
    
    # Inverse frequency weighting
    weights = 1.0 / (class_counts + 1e-6)
    
    # Extra emphasis on disease detection (non-healthy classes)
    for idx in range(len(weights)):
        if idx != 1:  # Not healthy class
            weights[idx] *= 2.0
    
    # Normalize
    weights = weights / weights.sum() * len(weights)
    
    return torch.tensor(weights, dtype=torch.float32)

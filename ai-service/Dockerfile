# Use Python 3.12 slim image for smaller size
FROM python:3.12-slim

# Set working directory
WORKDIR /app

# Set environment variables to prevent Python buffering
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

# Copy pyproject.toml first for better caching
COPY pyproject.toml .

# Install PyTorch CPU-only version first (much smaller and faster)
RUN pip install --no-cache-dir torch torchvision --index-url https://download.pytorch.org/whl/cpu

# Install FastAPI and Uvicorn first
RUN pip install --no-cache-dir fastapi uvicorn[standard] python-multipart

# Install remaining Python dependencies from pyproject.toml
RUN pip install --no-cache-dir .

# Copy application code
COPY models.py .
COPY data_utils.py .
COPY inference.py .
COPY app.py .

# Copy model files
COPY outputs/ outputs/

# Create non-root user for security
RUN useradd -m -u 1000 appuser && chown -R appuser:appuser /app
USER appuser

# Expose port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD python -c "import requests; requests.get('http://localhost:8000/health')" || exit 1

# Run the application
CMD ["python", "-m", "uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]

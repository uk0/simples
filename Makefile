# Makefile

# 声明伪目标，避免与同名文件冲突
.PHONY: build install clean

# 1) 构建目标
build:
	@echo "[INFO] Installing build tools..."
	rm -rf build simples.egg-info dist
	pip install --upgrade wheel setuptools
	@echo "[INFO] Building wheel package..."
	python setup.py bdist_wheel

# 2) 安装目标 (可选)
install:
	@echo "[INFO] Installing built wheel..."
	pip install dist/*.whl --force-reinstall

# 3) 清理目标 (可选)
clean:
	@echo "[INFO] Removing build artifacts..."
	rm -rf build dist *.egg-info
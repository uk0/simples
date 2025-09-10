# setup.py
from setuptools import setup, find_packages

setup(
    name="simples",
    version="0.0.1",
    packages=find_packages(),
    install_requires=[
        "markdown2==2.5.4",
        "markdownify==1.2.0",
        "pyqt6",
        "pathlib",
    ],
    entry_points={
        'console_scripts': [
            'notespush=watch_notes_service.notes_export:main',
        ],
    },
    python_requires=">=3.10",
)

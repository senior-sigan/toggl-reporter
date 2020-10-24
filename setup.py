import io
import os

import setuptools


def get_long_description():
    base_dir = os.path.abspath(os.path.dirname(__file__))
    file_name = os.path.join(base_dir, 'README.md')
    with io.open(file_name, encoding='utf-8') as fd:
        return fd.read()


def get_requirements():
    base_dir = os.path.abspath(os.path.dirname(__file__))
    file_name = os.path.join(base_dir, 'requirements.txt')
    with io.open(file_name, encoding='utf-8') as fd:
        li_req = fd.read().split('\n')
        return [req.strip() for req in li_req if len(req)]


setuptools.setup(
    name='toggl-reporter',
    version='0.0.2',
    author='Ilya Siganov',
    author_email='ilya.blan4@gmail.com',
    description='A tool to build daily reports from Toggle with nice grouping',
    long_description=get_long_description(),
    long_description_content_type='text/markdown',
    project_urls={
        'source': 'https://github.com/senior-sigan/toggl-reporter',  # noqa: E501
    },
    packages=setuptools.find_packages(),
    install_requires=get_requirements(),
    include_package_data=True,
    data_files=[
        ('requirements', ['requirements.txt']),
    ],
    zip_safe=False,
    entry_points={
        'console_scripts': [
            'toggl_reporter = toggl_reporter.__main__:main',
        ],
    },
    python_requires='>=3.6',
    license='MIT',
    classifiers=[
        'License :: OSI Approved :: MIT License',
        'Programming Language :: Python',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7',
        'Programming Language :: Python :: 3.8',
        'Environment :: Console',
        'Topic :: Utilities',
    ],
)

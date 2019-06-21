This example application showcases a complete CyberGRX export to an Excel file using the bulk-api-connector.  This application is done using Python, you should run all commands from this directory.

# Running the example
The first step is to configure a virtual environment for the application dependencies.  Depending on the version of Python that you are using the following commands will slightly differ.
- Python 2: `pip install virtualenv && virtualenv env`
- Python 3: `pip3 install virtualenv && python -m venv env`
- `source env/bin/activate`
- `pip install -r requirements.txt`
- `AUTH_TOKEN="${AUTH_TOKEN}" python export.py`
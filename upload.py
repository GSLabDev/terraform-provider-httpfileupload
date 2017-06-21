import os

from flask import Flask, request
from werkzeug.utils import secure_filename
from werkzeug.routing import BaseConverter


app = Flask(__name__)

UPLOAD_FOLDER = "/home/ismail/Documents/Uploads"
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER


@app.route('/upload', methods=['GET', 'POST', 'DELETE'])
def upload():
    if request.method == 'GET':
        return "The flask server is up and running"
    if request.method == 'POST':
        file = request.files['file']
        filename = secure_filename(file.filename)
        file.save(os.path.join(app.config['UPLOAD_FOLDER'], filename))
        app.config['UPLOADED_FILE'] = os.path.join(app.config['UPLOAD_FOLDER'], filename)
        return 'file uploaded successfully'

    if request.method == 'DELETE':
        file = request.files['file']
        os.remove(UPLOAD_FOLDER + "/" + file.filename)
        return 'file deleted successfully'


class RegexConverter(BaseConverter):
    def __init__(self, url_map, *items):
        super(RegexConverter, self).__init__(url_map)
        self.regex = items[0]


app.url_map.converters['regex'] = RegexConverter


@app.route('/upload/<regex("[\w./]{4,100}"):file_name>/', methods=['GET'])
def check(file_name):
    try:
        exit_status = os.path.exists(file_name)
        if exit_status != 0:
            error_message = 'File Does not exist'
            raise Exception(error_message)

    except Exception as e:
        return 0


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080, debug=True)

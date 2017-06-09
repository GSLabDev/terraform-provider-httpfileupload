import os

from flask import Flask, request
from werkzeug.utils import secure_filename

app = Flask(__name__)

UPLOAD_FOLDER = "/home/ismail/Documents/Uploads"
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER


@app.route('/upload', methods=['GET', 'POST', 'DELETE'])
def upload():
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


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080, debug=True)
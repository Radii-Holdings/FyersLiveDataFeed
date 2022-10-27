#this file will take the dataset created in raw to processing after specified interval
import schedule
import shutil
from datetime import date

today = date.today()
sourceFileName = today.strftime("%Y-%m-%d") + ".csv"
print('today file to operate ', sourceFileName)
#Dataset Path 
RAW_FILE_PATH = "./dataset/raw"
PROCESSING_FILE_PATH = './dataset/processing'
TIME_DELAY = 30 # MINUTES DATASET 
def CopyDataset():
    try:
        dest = shutil.copy(RAW_FILE_PATH+"/"+sourceFileName, PROCESSING_FILE_PATH+"/"+sourceFileName )
        print(dest)
    except shutil.SameFileError :
        print ('error while copying the dataset')
while True :
    print('forever')

schedule.every(TIME_DELAY).minutes.do(CopyDataset)


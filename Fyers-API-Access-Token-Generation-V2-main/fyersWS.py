from fyers_api import  websocket as ws
client_id="85VLN1I8IV-100"
token=open('fyers_access_token.txt','r').read().strip()
access_token = client_id+":"+token
data_type = 'symbolData'
symbol = ['NSE:TATAMOTORS-EQ','NSE:ITC-EQ']
def custom_message(ticks):
    print('ticks')
    print(ticks.response)

ws.FyersSocket.websocket_data = custom_message
fs = ws.FyersSocket(access_token, data_type,symbol)
fs.subscribe()
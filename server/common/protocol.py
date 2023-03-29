import logging
import struct

from common.utils import Bet
from common.sock import recvall, sendall

SUCCESS = 0
FAIL = 1
MORE_BATCHES = 1234
NO_MORE_BATCHES = 1235
CLOSE_CONECTION = -1

def _unpack_string(bytes, current_length):
    # Unpack the length of the string as a uint8
    length = struct.unpack('>B', bytes[current_length:current_length+1])[0]
    current_length += 1
    
    # Read the string data
    string_data = bytes[current_length:current_length+length].decode('utf-8')
    current_length += length
    
    # Return the string and the updated current length
    return string_data, current_length

def _deserialize_bet(bytes):
    # Initialize the current length variable to 0
    current_length = 0
    
    # Unpack the string fields using the unpack_string function
    name, current_length = _unpack_string(bytes, current_length)
    last_name, current_length = _unpack_string(bytes, current_length)
    id, current_length = _unpack_string(bytes, current_length)
    birthdate, current_length = _unpack_string(bytes, current_length)
    
    # Unpack the uint16 fields
    number = struct.unpack('>H', bytes[current_length:current_length+2])[0]
    current_length += 2
    
    agency_id = struct.unpack('>H', bytes[current_length:current_length+2])[0]
    
    # Create and return the Bet object
    return Bet(agency_id, name, last_name, id, birthdate, number)

def receive_bet(sock):
    answer = recvall(sock, 2)
    bet_size = struct.unpack('>H', answer)[0]
    if bet_size == MORE_BATCHES or bet_size == NO_MORE_BATCHES:
        return  CLOSE_CONECTION
    bet_byte_array = recvall(sock, int(bet_size))
    return _deserialize_bet(bet_byte_array)

def send_success(sock):
    sendall(sock, bytearray(SUCCESS.to_bytes(1, byteorder='big')), 1)

def send_fail(sock):
    sendall(sock, bytearray(FAIL.to_bytes(1, byteorder='big')), 1)


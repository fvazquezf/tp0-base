import socket
import logging
from common.protocol import receive_bet, send_fail, send_success
from common.utils import store_bets


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # the server
        while True:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except Exception as e:
                logging.error("action: accept_new_connection | result: fail ")
                break




    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            bet =  receive_bet(client_sock)
            addr = client_sock.getpeername()
            logging.info(f'action: receive_bet | result: success | ip: {addr[0]} | msg: {bet.document}')
            try:
                store_bets([bet])
                logging.info(f'action: store_bets | result: success | dni: {bet.document} | number: {bet.number}')
                send_success(client_sock)
            except BaseException as e:
                send_fail(client_sock)
                logging.error(f'action: store_bets | result: failed | error: {e}')


        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def _sigterm_handler(self, _signo, _stack_frame):
        logging.info(f'action: sigterm_handler | result: in_progress')
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        logging.info(f'action: sigterm_handler | result: success')


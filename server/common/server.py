import socket
import logging
from common.protocol import receive_bet, send_fail, send_not_finnished, send_success, send_winners, CLOSE_CONECTION, CLIENT_FINISHED, CLIENT_ASKING_RESPONSE
from common.utils import store_bets, find_winners

CLIENTS_AMMOUNT = 5

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
        finished_clients = set()
        winners = []
        # the server
        while True:
            try:
                client_sock = self.__accept_new_connection()
                logging.info(f'action: finished clients | {finished_clients}')
                if len(finished_clients) == CLIENTS_AMMOUNT:
                    if not winners:
                        logging.info(f'action: finding winners | result: in_progress ')
                        winners = find_winners()
                        logging.info(f'action: finding winners | result: success | Found: {len(winners)} ')
                    send_winners(client_sock, winners)
                else:
                    self.__handle_client_connection(client_sock, finished_clients)
            except Exception as e:
                logging.error(f"action: accept_new_connection | result: fail {e}")
                break




    def __handle_client_connection(self, client_sock, finished_clients):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try: 
            bets = []
            bet = receive_bet(client_sock)
            while bet != CLOSE_CONECTION and bet != CLIENT_FINISHED:
                bets.append(bet)
                bet = receive_bet(client_sock)
            logging.info(f'action: recieved 1 full batch | result: success')
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            client_sock.close()
            return
        if bet == CLIENT_ASKING_RESPONSE:
            send_not_finnished(client_sock)
            return
        try:
            store_bets(bets)
            logging.info(f'action: store_bets | result: success ')
            send_success(client_sock)
            logging.info(f'action: bet received | result: success | {bet}')
            if bet == CLIENT_FINISHED:
                finished_clients.add(bets[0].agency)
                logging.info(f'action: client finished | result: success | {bets[0].agency}')
                if len(finished_clients) == CLIENTS_AMMOUNT:
                    logging.info(f'action: sorteo | result: success')
        except BaseException as e:
            send_fail(client_sock)
            logging.error(f'action: store_bets | result: failed | error: {e}')
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


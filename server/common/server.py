import socket
import logging
import threading
from common.protocol import receive_bet, send_fail, send_success, send_winners, CLOSE_CONECTION, CLIENT_FINISHED, CLIENT_ASKING_RESPONSE
from common.utils import SafeBetStore
from common.safeList import SafeList

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
        winners = SafeList(CLIENTS_AMMOUNT) 
        betStorer = SafeBetStore()
        threads = {}
        # the server
        while True:
            try:
                client_sock = self.__accept_new_connection()
                thread = threading.Thread(target=self.__handle_client_connection, args=(client_sock, winners, betStorer))
                thread.start()
                threads[thread.ident] = thread # cuidado con el movimiento de los contexts del thread
                self.cleanThreads(threads)
            except Exception as e:
                logging.error(f"action: accept_new_connection | result: fail {e}")
                self.cleanThreads(threads)
                break

    def cleanThreads(self, threads):
        dead_threads = []
        for thread_id, thread in threads.items():
            if not thread.is_alive():
                dead_threads.append(thread_id)
                thread.join()
        for thread_id in dead_threads:
            del threads[thread_id]

    def __handle_client_connection(self, client_sock, winners, betStorer):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try: 
            bets = []
            bet = receive_bet(client_sock)
            while bet != CLOSE_CONECTION and bet != CLIENT_FINISHED and bet != CLIENT_ASKING_RESPONSE:
                bets.append(bet)
                bet = receive_bet(client_sock)
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            client_sock.close()
            return
        if bet == CLIENT_ASKING_RESPONSE:
            logging.info(f'action: client waiting for response {bet}')
            winners.wait_for_others()
            winners.fillList()
            send_winners(client_sock, winners)
            return
        try:
            betStorer.store_bets(bets)
            send_success(client_sock)
            if bet == CLIENT_FINISHED:
                logging.info(f'action: client finished | result: success | {bets[0].agency}')
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
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def _sigterm_handler(self, _signo, _stack_frame):
        logging.info(f'action: sigterm_handler | result: in_progress')
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        logging.info(f'action: sigterm_handler | result: success')

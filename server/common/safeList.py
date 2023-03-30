import logging
import threading

from common.utils import find_winners

class SafeList:
    def __init__(self, clientsAmmount):
        self.lock = threading.Lock()
        self.list = []
        self.barrier = threading.Barrier(clientsAmmount)

    def fillList(self):
        with self.lock:
            if not len(self.list):
                self.list.extend(find_winners())

    def filter_bets_by_agency(self, clientID):
        with self.lock:
            return [item for item in self.list if item.agency == clientID]

    def wait_for_others(self):
        self.barrier.wait()

    def __len__(self):
        with self.lock:
            return len(self.list)
        
    def copy(self):
        with self.lock:
            return list(self.list)
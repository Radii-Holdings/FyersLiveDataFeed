
########################


from __future__ import annotations
from abc import ABC, abstractmethod
from random import randrange
from typing import List


class Subject(ABC):
    """
    The Subject interface declares a set of methods for managing subscribers.
    """

    @abstractmethod
    def attach(self, observer: Observer) -> None:
        """
        Attach an observer to the subject.
        """
        pass

    @abstractmethod
    def detach(self, observer: Observer) -> None:
        """
        Detach an observer from the subject.
        """
        pass

    @abstractmethod
    def notify(self) -> None:
        """
        Notify all observers about an event.
        """
        pass


class BEACONSubject(Subject):
    """
    The Subject owns some important state and notifies observers when the state
    changes.
    """

    _state: int = None
    """
    For the sake of simplicity, the  Beacon Subject's state, is 
    """

    _observers: List[Observer] = []
    """
    List of subscribers. In real life, the list of subscribers can be stored
    more comprehensively (categorized by event type, etc.).
    """

    def attach(self, observer: Observer) -> None:
        print("Subject: Attached an observer.")
        self._observers.append(observer)

    def detach(self, observer: Observer) -> None:
        self._observers.remove(observer)

    """
    The subscription management methods.
    """

    def notify(self) -> None:
        """
        Trigger an update in each subscriber.
        """

        print("Subject: Notifying observers...")
        for observer in self._observers:
            observer.update(self)

    def some_business_logic(self) -> None:
        """
        Usually, the subscription logic is only a fraction of what a Subject can
        really do. Subjects commonly hold some important business logic, that
        triggers a notification method whenever something important is about to
        happen (or after it).
        """

        print("\nSubject: I'm doing something important.")
        self._state = randrange(0, 10)

        print(f"Subject: My state has just changed to: {self._state}")
        self.notify()


class Observer(ABC):
    """
    The Observer interface declares the update method, used by subjects.
    """

    @abstractmethod
    def update(self, subject: Subject) -> None:
        """
        Receive update from subject.
        """
        pass


"""
Concrete Observers react to the updates issued by the Subject they had been
attached to.
"""


class ConcreteObserverA(Observer):
    def update(self, subject: Subject) -> None:
        if subject._state < 3:
            print("ConcreteObserverA: Reacted to the event")


class ConcreteObserverB(Observer):
    def update(self, subject: Subject) -> None:
        if subject._state == 0 or subject._state >= 2:
            print("ConcreteObserverB: Reacted to the event")

# This portion of the code is repeatedly processomg tje Raw source and making log actions
# 
# this file will take the dataset created in raw to processing after specified interval
import schedule
import shutil
from datetime import date
import time

today = date.today()
sourceFileName = today.strftime("%Y-%m-%d") + ".csv"
print('today file to operate ', sourceFileName)
#Dataset Path 
RAW_FILE_PATH = "./dataset/raw"
PROCESSING_FILE_PATH = './dataset/processing'
NEW_DATA_OFFSET = 1 # MINUTES DATASET 
BEACON_OFFSET = 5 # SECONDS RECALCULATION LOOPS
AI_OFFSET =30 # SECONDS PREDICTORS

###########################################################################################

# dataset operative functions

################################################################################

def copy_dataset():
    try:
        dest = shutil.copy(RAW_FILE_PATH+"/"+sourceFileName, PROCESSING_FILE_PATH+"/"+sourceFileName )
        print(dest)
        return dest
    except shutil.SameFileError :
        print ('error while copying the dataset')

def BEACON_observers():
    '''
        The BEACON observers shall feed upon market parameters and technical indicators / 
        be called on specified intervals
    '''
    subject = BEACONSubject()

    observer_a = ConcreteObserverA()
    subject.attach(observer_a)

    observer_b = ConcreteObserverB()
    subject.attach(observer_b)

    subject.some_business_logic()
    subject.some_business_logic()

    subject.detach(observer_a)

    subject.some_business_logic()

#####################################################################################################

# Can opener entry point for the processing of realtime data

if __name__ == "__main__":
    # The client code.



    schedule.every(NEW_DATA_OFFSET).minutes.do(copy_dataset)

    # schedule.every(BEACON_OFFSET).minutes.do(BEACON_observers)

    while True :
        schedule.run_pending()
        time.sleep(1)
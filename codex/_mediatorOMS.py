from __future__ import annotations
from abc import ABC

BUY_CNC = 'buy'
SELL_CNC = 'sell'

class Mediator(ABC):
    """
    The Mediator interface declares a method used by components to notify the
    mediator about various events. The Mediator may react to these events and
    pass the execution to other components.
    """

    def notify(self, sender: object, event: str) -> None:
        pass


class ConcreteMediator(Mediator):
    '''
    TODO add all the algorthms inside this mediator
    '''
    def __init__(self, component1: SpikeOMSEvents) -> None:
        self._component1 = component1
        self._component1.mediator = self
        self._component2 = _SellEvent
        self._component2.mediator = self

    def notify(self, sender: object, event: str) -> None:
        if event == BUY_CNC:
            print("Mediator reacts on A and triggers following operations:")
            self._component2.do_c()
        elif event == "D":
            print("Mediator reacts on D and triggers following operations:")
            self._component1.do_b()
            self._component2.do_c()


class BaseComponent:
    """
    The Base Component provides the basic functionality of storing a mediator's
    instance inside component objects.
    """

    def __init__(self, mediator: Mediator = None) -> None:
        self._mediator = mediator

    @property
    def mediator(self) -> Mediator:
        return self._mediator

    @mediator.setter
    def mediator(self, mediator: Mediator) -> None:
        self._mediator = mediator


"""
Concrete Components implement various functionality. They don't depend on other
components. They also don't depend on any concrete mediator classes.
"""


class SpikeOMSEvents(BaseComponent):
    def notify_buy(self) -> None:
        print(" Spike event -- new buy arrived")
        self.mediator.notify(self,sender= 24, event= BUY_CNC)

    def do_buy(self) -> int:
        print("Component 1 does B.")
        self.mediator.notify(self, "B")
        return 0

#internal algorithmic classes
class _SellEvent(BaseComponent):
    def do_c(self) -> None:
        print("Spike event -- new buy arrived")
        self.mediator.notify(self, "C")

    def do_d(self) -> None:
        print("Component 2 does D.")
        self.mediator.notify(self, "D")


def placeOrder(json:object):
    # The client code.
    c1 = SpikeOMSEvents()
    ConcreteMediator(c1)

    print("Mediator range call for processing buy and sell order together ")
    c1.notify_buy()

    print("\n", end="")

    
    
from trustix_dash.lib.defer import DeferStack
import itertools
import typing


T = typing.TypeVar("T")


def flatten(seq: typing.Iterable[typing.Iterable[T]]) -> typing.Iterator[T]:
    return itertools.chain.from_iterable(seq)


def unique(seq: typing.Iterable[T]) -> typing.Generator[T, None, None]:
    seen: typing.Set[T] = set()
    for i in seq:
        if i not in seen:
            seen.add(i)
            yield i


__all__ = ("DeferStack", "flatten")

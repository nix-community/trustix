# Trustix
# Copyright (C) 2021 Tweag IO

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

from tortoise.models import Model
from typing import Union, Type
from tortoise import fields
import base64


__all__ = ("BinaryField",)


class BinaryField(fields.CharField):
    """
    An implementation of a binary field using a CharField and base85 encoding

    BLOB's are not indexable in all databases, hence we use VarChar
    """

    @staticmethod
    def encode_value(value: Union[str, bytes]) -> str:
        if isinstance(value, str):
            return value
        return base64.b85encode(value).decode()

    def to_db_value(
        self, value: Union[str, bytes], instance: Union[Type[Model], Model]
    ) -> str:
        return self.encode_value(value)

    def to_python_value(self, value: Union[bytes, str]) -> bytes:
        if isinstance(value, bytes):
            return value
        return base64.b85decode(value)

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

from tortoise import models
from tortoise import fields

from trustix_nix_reprod.models.meta import BaseMeta


class Evaluation(models.Model):
    commit = fields.CharField(
        max_length=40,
        index=True,
        pk=True,
    )

    timestamp = fields.DatetimeField(auto_now_add=True)

    def __str__(self):
        return self.commit

    class Meta(BaseMeta):
        pass

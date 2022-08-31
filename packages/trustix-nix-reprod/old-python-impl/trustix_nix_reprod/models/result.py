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

from trustix_nix_reprod.models import fields as trustix_fields
from trustix_nix_reprod.models.meta import BaseMeta, app_name


class DerivationOutputResult(models.Model):
    output = fields.ForeignKeyField(
        f"{app_name}.DerivationOutput",
        to_field="input_hash",
        on_delete=fields.CASCADE,
        db_constraint=False,
        index=True,
    )
    # TODO: Turn into indexed CharField
    output_hash = trustix_fields.BinaryField(max_length=40)
    log = fields.ForeignKeyField(
        f"{app_name}.Log",
        on_delete=fields.CASCADE,
    )

    # def __str__(self):
    #     return f"{self.output_id}"
    #     # return f"{self.log_id}({self.output.derivation_id, self.output.output})"

    class Meta(BaseMeta):
        unique_together = (("output", "log"),)

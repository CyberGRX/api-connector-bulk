#########################################################################
#    _________        ___.                   ______________________  ___
#    \_   ___ \___.__.\_ |__   ___________  /  _____/\______   \   \/  /
#    /    \  \<   |  | | __ \_/ __ \_  __ \/   \  ___ |       _/\     /
#    \     \___\___  | | \_\ \  ___/|  | \/\    \_\  \|    |   \/     \
#     \______  / ____| |___  /\___  >__|    \______  /|____|_  /___/\  \
#            \/\/          \/     \/               \/        \/      \_/
#
#

import os
import json
import requests
from openpyxl import Workbook
from tqdm import tqdm
from glom import glom

def sheet_writer(wb, name, columns, mapping):
    def builder(sheet):
        for idx, injector in enumerate(columns):
            sheet.cell(row=1, column=1+idx).value = injector[0]

        row = 2
        def writer(blob):
            nonlocal row
            transformed = glom(blob, mapping)
            for idx, injector in enumerate(columns):
                value = transformed[injector[1]]
                if value is None:
                    continue
                    
                sheet.cell(row=row, column=1+idx).value = transformed[injector[1]]
            row += 1

        return writer

    return builder(wb[name])
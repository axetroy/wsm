#!/usr/bin/python
# -*- coding: UTF-8 -*-

# Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.

import os


def walk(root_dir):
    for root, dirs, files in os.walk(root_dir, topdown=False):
        for filename in files:
            if filename == "node_modules" or filename == "vendor":
                continue

            if os.path.splitext(filename)[-1] == ".go":
                filepath = os.path.join(root, filename)

                with open(filepath, "r") as reader:
                    originContent = reader.read()
                    reader.close()
                    if not originContent.startswith("// Copyright"):
                        with open(filepath, "w") as writer:
                            writer.write("// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.\n")
                            writer.write(originContent)

            elif os.path.splitext(filename)[-1] == ".js":
                filepath = os.path.join(root, filename)

                with open(filepath, "r") as reader:
                    originContent = reader.read()
                    reader.close()
                    if not originContent.startswith("// Copyright"):
                        with open(filepath, "w") as writer:
                            writer.write("// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.\n")
                            writer.write(originContent)

            elif os.path.splitext(filename)[-1] == ".vue":
                filepath = os.path.join(root, filename)

                with open(filepath, "r") as reader:
                    originContent = reader.read()
                    reader.close()
                    if not originContent.startswith("<!-- Copyright"):
                        with open(filepath, "w") as writer:
                            writer.write("<!-- Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0. "
                                         "-->\n")
                            writer.write(originContent)


walk("./internal")
walk("./cmd")
walk("./frontend")

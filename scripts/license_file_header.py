#!/usr/bin/python
# -*- coding: UTF-8 -*-

# Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.

import os

exclude = {"node_modules", "vendor", ".nuxt", "scripts", "logs", "docker"}


def walk(root_dir):
    for root, dirs, files in os.walk(root_dir, topdown=True):
        dirs[:] = [d for d in dirs if d not in exclude]
        for filename in files:
            filepath = os.path.join(root, filename)

            if os.path.splitext(filename)[-1] == ".go":
                with open(filepath, "r") as reader:
                    originContent = reader.read()
                    reader.close()
                    if not originContent.startswith("// Copyright"):
                        with open(filepath, "w") as writer:
                            writer.write("// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.\n")
                            writer.write(originContent)

            elif os.path.splitext(filename)[-1] == ".js":
                with open(filepath, "r") as reader:
                    originContent = reader.read()
                    reader.close()
                    if not originContent.startswith("// Copyright"):
                        with open(filepath, "w") as writer:
                            writer.write("// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.\n")
                            writer.write(originContent)

            elif os.path.splitext(filename)[-1] == ".vue":
                with open(filepath, "r") as reader:
                    originContent = reader.read()
                    reader.close()
                    if not originContent.startswith("<!-- Copyright"):
                        with open(filepath, "w") as writer:
                            writer.write("<!-- Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0. "
                                        "-->\n")
                            writer.write(originContent)


walk(".")

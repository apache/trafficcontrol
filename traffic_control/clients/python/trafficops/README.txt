- Introduction

This directory contains the Traffic Ops Python Client for Python 2.x and Python 3.x.

NOTE: This client has only been tested against Python 2.7 and Python 3.6.  Other versions
      may work, so, you mileage may vary.
      
NOTICE: As of Apache Traffic Control release version 4.0.0, Python 2 will no longer be supported! Developers are encouraged
        to switch to version 3 as soon as possible. The CPython project has a free tool for automating much of the process of
        translating Python2 code to Python3, read more about it here: https://docs.python.org/3/library/2to3.html

- Installation

NOTE: Assuming in a already newly created and activated virtualenv

From github.com via pip:

# pip install git+https://github.com/apache/trafficcontrol.git#"egg=trafficops&subdirectory=traffic_control/clients/python/trafficops"

or

# pip install git+ssh://git@github.com/apache/trafficcontrol.git#"egg=trafficops&subdirectory=traffic_control/clients/python/trafficops"

Local Installation from cloned 'trafficcontrol' repository:

NOTE: 'trafficcontrol' => https://www.github.com/apache/trafficcontrol

NOTE: Assuming cd'd to the root of a forked/cloned 'trafficcontrol' repository

    1. Activate your virtualenv environment (Python Virtual Environment)
        E.g.
            NOTE: Where <virtual environment directory> is, for example, '~/VENV/my_venv'.
            $ source ~/VENV/my_venv/bin/activate
            (my_venv) $
    2. Install the software for the correct environment in the activated virtualenv
        (my_venv) $ cd <root of 'trafficcontrol' repository>/traffic_control/clients/python/trafficops
        (my_venv) $ python setup.py install
        ...
        (my_venv) $ cd <root of 'trafficcontrol' repository>
    3. Test Package is installed correctly
        (my_venv) $ python
        (my_venv) $ Python 3.6.1 (default, Apr  4 2017, 09:40:21)
                    [GCC 4.2.1 Compatible Apple LLVM 8.1.0 (clang-802.0.38)] on darwin
                    Type "help", "copyright", "credits" or "license" for more information.
                    >>> import trafficops
                    >>> dir(trafficops)
                    ['LoginError', 'OperationError', 'RestApiSession', 'TOSession', '__builtins__', '__cached__',
                     '__doc__', '__file__', '__loader__', '__name__', '__package__', '__path__', '__spec__',
                      '__version__', 'api_request', 'default_headers', 'restapi', 'tosession']
                    >>> tos = trafficops.TOSession(host_ip=u'to.somedomain.net', verify_cert=True)
                    >>> tos.login(u'someuser', u'someuser123')
                    >>> cdns, response = tos.get_cdns()
                    >>> exit()
        (my_venv) $

        NOTE: No errors/exceptions (warnings are okay) means you should be good to go.

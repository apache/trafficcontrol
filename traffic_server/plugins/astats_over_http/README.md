Compile:
  tsxs -c astats_over_http.c -o astats_over_http.so
Install:
  sudo tsxs -o astats_over_http.so -i

Add to the plugin.conf:

  astats_over_http.so path=${path}

start traffic server and visit http://[ip]:[port]/${path}

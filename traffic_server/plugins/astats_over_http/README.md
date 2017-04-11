Compile:
  tsxs -c astats_over_http.c -o astats_over_http.so
Install:
  sudo tsxs -o astats_over_http.so -i

Add to the plugin.conf:

  astats_over_http.so path=${path}

start traffic server and visit http://[ip]:[port]/${path}

Rpm Builds

  Two spec files are provided.  astats_over_http.spec requires a tar ball of this directoy 
  named astats_over_htt-.tar.gz is copied to the rpmbuild/SOURCES directory.  The second
  astats-git-build, checks out the source from the git repo and builds the rpm.

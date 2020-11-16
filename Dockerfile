FROM scratch
COPY /css /css
COPY /html html
COPY /js js
COPY /fonts fonts
COPY /linux /
CMD ["/system_webservice_beta"]
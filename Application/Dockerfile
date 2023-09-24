FROM golang


RUN apt-get install -y curl

# Set the Current Working Directory inside the container
WORKDIR /app

#COPY app
COPY app/priceapp .


COPY . . 

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the program 
CMD ["./priceapp"]

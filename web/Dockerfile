#FROM node:alpine AS dev
#COPY ./sapper .
#EXPOSE 3000
#CMD ["npm", "run", "start"]

FROM node:alpine
WORKDIR /goTemp/web/sapper
COPY ./sapper/package*.json ./
RUN npm install
COPY ./sapper .
EXPOSE 3000
#EXPOSE 10000
ENV HOST=0.0.0.0
CMD [ "npm", "build" ]
CMD [ "npm", "start" ]
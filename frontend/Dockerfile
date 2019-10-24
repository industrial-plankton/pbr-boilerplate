FROM node

# if left blank app will run with dev settings
# to build production image run:
# $ docker build ./frontend --build-args app_env=production
ENV NPM_CONFIG_LOGLEVEL warn

# This is how you load those "args" values in the docker-compose.yml file
# into the container's local environment
ARG APP_ENV
ENV NODE_ENV $APP_ENV

RUN mkdir -p /frontend
WORKDIR /frontend
COPY ./ ./

RUN npm i

CMD if [ ${NODE_ENV} = production ]; \
	then \
	yarn run start:production; \
	else \
	yarn start; \
	fi

EXPOSE 3000

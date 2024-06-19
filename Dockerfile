FROM public.ecr.aws/lambda/go:latest

# Copy env file
COPY .env .env

# Copy function code
COPY bootstrap ${LAMBDA_TASK_ROOT}

# Set the CMD to your handler (could also be done as a parameter override outside of the Dockerfile)
CMD [ "bootstrap" ]
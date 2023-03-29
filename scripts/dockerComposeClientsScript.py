import sys

args = sys.argv

if len(args) != 2:
    print("Please provide only 1 input parameter.")
    exit(-1)

try:
    input_param = int(args[1])
except:
    print("Only integers are accepted")
    exit(-1)

file = open("../docker-compose-dev.yaml", "w")

file.write("version: '3.9'\n"
           "name: tp0\n"
           "services:\n"
           "  server:\n"
           "    container_name: server\n"
           "    image: server:latest\n"
           "    entrypoint: python3 /main.py\n"
           "    environment:\n"
           "      - PYTHONUNBUFFERED=1\n"
           "      - LOGGING_LEVEL=DEBUG\n"
           "    volumes:\n"
           "      - ./server/config.ini:/config.ini\n"
           "    networks:\n"
           "      - testing_net\n\n")

for i in range (input_param):
    client = "client"+str(i+1)
    file.write(
            "  "+client+":\n"
            "    container_name: "+client+"\n"
            "    image: client:latest\n"
            "    entrypoint: /client\n"
            "    environment:\n"
            "      - CLI_AGENCY="+str(i+1)+"\n"
            "      - CLI_LOG_LEVEL=DEBUG\n"
            "      - CLI_BATCHMAXSIZE=10000\n"
            "    networks:\n"
            "      - testing_net\n"
            "    volumes:\n"
            "      - ./client/config.yaml:/config.yaml\n"
            "      - ./.data/dataset/agency-"+str(i+1)+".csv:/dataset.csv\n"
            "    depends_on:\n"
            "      - server\n\n"
        )

file.write(
        "networks:\n"
        "  testing_net:\n"
        "    ipam:\n"
        "      driver: default\n"
        "      config:\n"
        "        - subnet: 172.25.125.0/24\n"
    )

file.close()
exit(0)
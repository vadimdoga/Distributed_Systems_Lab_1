import pika
import json

credentials = pika.PlainCredentials('guest', 'guest')
connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost', port=5672, virtual_host='/', credentials=credentials))
channel = connection.channel()

channel.queue_declare(queue='ORDER_CREATED', durable=True)

channel.basic_publish(exchange='',
                      routing_key='ORDER_CREATED',
                      body=json.dumps({
                          "id": 23,
                          "magic": "exists"
                      }))
print(" [x] Sent json")

connection.close()

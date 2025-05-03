import pygame as pg
import socket
import pickle

GREEN = (0, 255, 0)
WHITE = (255, 255, 255)
BLACK = (0, 0, 0)
RED = (255, 0, 0)

HOSTNAME = "localhost"
PORT = 8080

pg.init()
screen = pg.display.set_mode((800, 600))
screen.fill(WHITE)
pg.display.flip()

pg.display.set_caption("Paint Server")

font = pg.font.SysFont("Comic Sans MS", 30)
text_surface = font.render("Server", True, RED)
screen.blit(text_surface, (screen.get_width() - text_surface.get_width() - 10, 10))

server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_socket.bind((HOSTNAME, PORT))
server_socket.listen(1)
client_socket, address = server_socket.accept()

running = True
while running:
    for event in pg.event.get():
        if event.type == pg.QUIT:
            running = False

    try:
        data = client_socket.recv(4096)
        if not data:
            break

        start, end = pickle.loads(data)
        pg.draw.line(screen, GREEN, start, end, 2)
    except (ConnectionResetError, pickle.UnpicklingError):
        break

    screen.blit(text_surface, (screen.get_width() - text_surface.get_width() - 10, 10))

    pg.display.flip()

client_socket.close()
server_socket.close()
pg.quit()

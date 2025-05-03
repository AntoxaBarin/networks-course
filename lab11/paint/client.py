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

pg.display.set_caption("Paint Client")

font = pg.font.SysFont("Comic Sans MS", 30)
text_surface = font.render("Client", True, RED)
screen.blit(text_surface, (screen.get_width() - text_surface.get_width() - 10, 10))

client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client_socket.connect((HOSTNAME, PORT))

drawing = False
last = None


while True:
    for event in pg.event.get():
        if event.type == pg.MOUSEBUTTONDOWN:
            drawing = True
        elif event.type == pg.MOUSEBUTTONUP:
            drawing = False
            last = None

    if drawing:
        mouse_pos = pg.mouse.get_pos()
        if last != None:
            pg.draw.line(screen, GREEN, last, mouse_pos, 2)
            client_socket.send(pickle.dumps((last, mouse_pos)))
        last = mouse_pos
    pg.display.flip()

    screen.blit(text_surface, (screen.get_width() - text_surface.get_width() - 10, 10))

pg.quit()

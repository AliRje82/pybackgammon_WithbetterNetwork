import pygame
from game import Board, Dice, Piece, OtherMouse
import json



class Game():
    def __init__(self, connection):
        self._running = True
        self._screen = None
        self.reset_sound = None
        self.run_server = True
        self.connection = connection
        self.size = self.width, self.height = 1000, 500
        self.board = Board(self)
        self.dice = Dice(self)
        self.init_pieces()
        self.player_count = 0
        self.other_mouse = OtherMouse()

    def init_pieces(self, send=True):
        self.pieces = list()
        self.fields = [[] for _ in range(24)]
        self.fields[0] = [True] * 2
        self.fields[5] = [False] * 5
        self.fields[7] = [False] * 3
        self.fields[11] = [True] * 5
        self.fields[23] = [False] * 2
        self.fields[18] = [True] * 5
        self.fields[16] = [True] * 3
        self.fields[12] = [False] * 5
        self.pieces = list()
        self.piece_size = 42
        self.ping_iter = 0
        ident = 1
        for field_id, field in enumerate(self.fields):
            top = field_id // 12 == 1
            relative_field_id = field_id % 12
            for piece_id, is_black in enumerate(field):
                offset_x = self.board.bounding_box_width + self.board.triangle_width//2 + \
                    self.board.triangle_width * relative_field_id + \
                    (relative_field_id // 6) * self.board.offset_x
                x = offset_x if top else self.width - offset_x
                (relative_field_id // 6) * self.board.offset_x
                y = self.board.bounding_box_width + self.piece_size * \
                    (piece_id*2+1) if top else self.height - \
                    self.piece_size * (piece_id*2+1) - \
                    self.board.bounding_box_width
                self.pieces.append(
                    Piece(self, ident=ident, pos=(x, y), black=is_black))
                ident += 1
        self.dice.reset()

        if self.reset_sound is not None:
            self.reset_sound.play()
            if send:
                message = json.dumps({"action": "resetboard"}) 
                self.connection.send(message.encode())

    def send_gamestate(self):
        for p in self.pieces:
            p.send_move()
        self.dice.send_state()
        self.dice.send_eyes()

    def on_init(self):
        pygame.init()
        pygame.mixer.init()
        self.reset_sound = pygame.mixer.Sound('assets/sound/button.wav')
        self.impact_sound = pygame.mixer.Sound('assets/sound/impact.wav')
        self.font = pygame.font.Font(pygame.font.get_default_font(), 22)
        pygame.display.set_caption('Backgammon')
        self.clock = pygame.time.Clock()
        self._screen = pygame.display.set_mode(
            self.size, pygame.HWSURFACE | pygame.DOUBLEBUF)
        self._running = True

    def ping(self):
        message = json.dumps({"action": "ping"}) 
        self.connection.send(message.encode())

    def keep_connection_alive(self):
        self.ping_iter = (self.ping_iter + 1) % 240
        if self.ping_iter == 0:
            self.ping()

    def on_event(self, event):
        if event.type == pygame.QUIT:
            self._running = False
        elif event.type == pygame.KEYDOWN:
            if event.key == pygame.K_SPACE:
                self.dice.roll()
            elif event.key == pygame.K_ESCAPE:
                self.init_pieces()
        else:
            self.handle_piece_events(event)
            if event.type == pygame.MOUSEMOTION:
                message = json.dumps({'action': 'mousemotion', 'pos': event.pos}) 
                self.connection.send(message.encode())

    def handle_piece_events(self, event):
        for idx, piece in enumerate(self.pieces):
            if piece.handle_event(event):
                if idx == 0:
                    return
                for idx2, piece2 in enumerate(self.pieces):
                    if idx == idx2:
                        continue
                    if piece.rect.colliderect(piece2.rect):
                        return
                else:
                    self.pieces.insert(0, self.pieces.pop(idx))
                return

        if self.dice.handle_event(event):
            return

        pygame.mouse.set_cursor(pygame.SYSTEM_CURSOR_ARROW)

    # def on_loop(self):
    #     self.keep_connection_alive()
    #     self.connection.Pump()
    #     self.Pump()
    #     if self.run_server:
    #         self.server.Pump()

    def on_render(self):
        self.board.render(self._screen)
        for piece in self.pieces[::-1]:
            piece.update(self._screen)
        self.dice.render(self._screen)
        self.other_mouse.render(self._screen)
        pygame.display.flip()

    def on_cleanup(self):
        pygame.quit()

    def on_execute(self):
        if self.on_init() == False:
            self._running = False

        while (self._running):
            self.clock.tick(60)
            for event in pygame.event.get():
                self.on_event(event)
            # self.on_loop()
            self.on_render()
        self.on_cleanup()

# A minimal wrapper over the normal Karel functions to make working with it easier

# unfournately race conditions are possible and only 1 Karel can exist at once
current_rotation = 0 # lets commit the cardinal sin of creating a mutatable global variable
def rotate(degree):
    global current_rotation
    for i in range(0, degree/90):
        turn_left()
    current_rotation += degree - (degree%90)
    current_rotation %= 360

def move_up(steps):
    rotate(90)
    for i in range(0, steps):
        move()
    rotate(270)

# basically just an alliss
def move_right(steps):
    for i in range(0, steps):
        move()
    
def move_left(steps):
    rotate(180)
    for i in range(0, steps):
        move()
    rotate(180)

def move_down(steps):
    rotate(270)
    for i in range(0, steps):
        move()
    rotate(90)
    
def spawn_ball(amount):
    for i in range(0, amount):
        put_ball() 
        
def eat_ball(amount):
    for i in range(0, amount):
        take_ball() 

def dump_ball(depth, amount):
    move_down(depth)
    spawn_ball(amount)
    move_up(depth)

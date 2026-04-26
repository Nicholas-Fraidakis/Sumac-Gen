# Sumac-Gen
Sumac-Gen is a small little Karel auto-solver and code generator!

It works by internally solving a Karel problem and recording the movements and objective/actions (picking up and placing a ball for example). Then does a small conversion to text (well Python code but it is basically just a text conversion). 

Also the outputed code relies on the functions from sumac-lib.py (which contains some helper functions to move Karel)

## Status
Right now sumac-gen is in a functional but very early and prototype-like state, so the codebase is pretty messy right now and needs to be cleaned up! 
Also gotta rename a couple internal names because this project because this started out as joke project (still somewhat is), but is now meant to be a polished (eventually) and usuable tool. Also now it is open source too on my github account that is for my more serious projects I have to make it a little more professional :)

That being said it can solve maps where:
- Pick up and place balls
- Karel has to navagate/pathfind (with or without walls) to an end position

Here is some things that need to be done before I am satified:
- Usabilty:
    - Add CLI Support instead of hardcoding Karel maps
    - Add Semi-Graphical Support (use robotgo to capture Karel maps on screen and parse them into usable representation) then output code to file
- Karel Functionally:
    - Add end rotation support
    - Add painting and other Karel actions
- Code gen:
    - Optimization:
        - Add ability to collapse code into a for loop
        - Add ability for repeated code to be turned into a function

## Credits
- Made by Nicholas Fraidakis in 2026 :D

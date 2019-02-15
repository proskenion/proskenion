#!/usr/bin/env python
# -*- coding: utf-8 -*-

from graphviz import Digraph

N = 10      # Creator数
M = 4       # Peer数

def cr(a):
    return "C"+str(a)

def pr(a):
    return "P"+str(a)

def initG():
    # formatはpngを指定(他にはPDF, PNG, SVGなどが指定可)
    G = Digraph(format='png')

    C = Digraph(name='creator', format='png')
    C.attr('node', shape='circle', style='filled', fillcolor='#EEEEEE')
    C.attr('graph', rank='source')

    P = Digraph(name='peer', format='png')
    P.attr('node', shape='trapezium', style='filled', fillcolor='#FFEECC')
    P.attr('graph', rank='sink')

    # ノードの追加
    for i in range(N):
        C.node(cr(i), cr(i))

    for i in range(M):
        if i == 0:
            P.node(pr(i), 'root')
        else:
            P.node(pr(i), pr(i))

    G.subgraph(C)
    G.subgraph(P)
    return G

G = initG()

for i in range(N):
    G.edge(cr(i), pr(i%M), color='#CCCCCC')
G.node(pr(0),'root',fillcolor='#FF99CC')

# binary_tree.pngで保存
G.render('_demo_graph01')

def firstEdge(G):
    # 辺の追加
    for i in range(N):
        G.edge(cr(i), cr((i+1)%N))
        if i < 8:
            G.edge(cr(i), cr((i+2)%N))
        if i < 6:
            G.edge(cr(i), cr((i+3)%N))
        if i < 4:
            G.edge(cr(i), cr((i+4)%N))

        if i > 5:
            G.node(cr(i), label=cr(i), style='filled', fillcolor='#FFCC33')

G = initG()
firstEdge(G)

for i in range(N):
    if i > 5:
        G.edge(cr(i),pr(i%M), color='#FF1122')
    else:
        G.edge(cr(i), pr(i%M), color='#CCCCCC')

for i in range(M):
    G.node(pr(i) ,fillcolor='#FF99CC')

G.render('_demo_graph02')

def secondEdge(G):
    # 辺の追加
    for i in range(N):
        if i < 5:
            G.edge(cr(i), cr(2))
        else:
            G.edge(cr(i), cr(4))
        G.node(cr(i), cr(i), fillcolor='#EEEEEE')

G = initG()
firstEdge(G)
secondEdge(G)

for i in range(N):
    if i in [0,1,3,4,5,9]:
        G.edge(cr(i), pr(i%M), color='#CCCCCC')

G.node(cr(2),cr(2),fillcolor='#FFCC33')
G.node(cr(4),cr(4),fillcolor='#FFCC33')
G.node(cr(7),cr(7),fillcolor='#FFCC33')
G.node(cr(8),cr(8),fillcolor='#FFCC33')
G.edge(cr(2),pr(2),color='#FF1122')
G.edge(cr(4),pr(2),color='#FF1122')
G.edge(cr(7),pr(3),color='#FF1122')
G.edge(cr(8),pr(0),color='#FF1122')

for i in range(M):
    if i != 1:
        G.node(pr(i) ,fillcolor='#FF99CC')

G.render(('_demo_graph03'))




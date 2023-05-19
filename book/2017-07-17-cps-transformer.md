---
layout: post
title: CPS 王垠 大佬的代码[测试过]
categories: lisp,cps
description: 转载
keywords: lisp, Script
---


发这个贴的原因，是因为 是因为 是因为啥来着，忘了 那就是因为你，因为你 。


## 摘要: `CPS.scm` ,`pmatch.scm`
   * `scheme`代码。


##  Stat

   * `CPS.scm` Code
   * `pmatch.scm` Code


* CPS

```scheme
;; A simple CPS transformer which does proper tail-call and does not
;; duplicate contexts for if-expressions.

;; author: Yin Wang (yw21@cs.indiana.edu)


(load "pmatch.scm")


(define cps
(lambda (exp)
(letrec
([trivial? (lambda (x) (memq x '(zero? add1 sub1)))]
[id (lambda (v) v)]
[ctx0 (lambda (v) `(k ,v))] ; tail context
[fv (let ([n -1])
(lambda ()
(set! n (+ 1 n))
(string->symbol (string-append "v" (number->string n)))))]
[cps1
(lambda (exp ctx)
(pmatch exp
[,x (guard (not (pair? x))) (ctx x)]
[(if ,test ,conseq ,alt)
(cps1 test
(lambda (t)
(cond
[(memq ctx (list ctx0 id))
`(if ,t ,(cps1 conseq ctx) ,(cps1 alt ctx))]
[else
(let ([u (fv)])
`(let ([k (lambda (,u) ,(ctx u))])
(if ,t ,(cps1 conseq ctx0) ,(cps1 alt ctx0))))])))]
[(lambda (,x) ,body)
(ctx `(lambda (,x k) ,(cps1 body ctx0)))]
[(,op ,a ,b)
(cps1 a (lambda (v1)
(cps1 b (lambda (v2)
(ctx `(,op ,v1 ,v2))))))]
[(,rator ,rand)
(cps1 rator
(lambda (r)
(cps1 rand
(lambda (d)
(cond
[(trivial? r) (ctx `(,r ,d))]
[(eq? ctx ctx0) `(,r ,d k)] ; tail call
[else
(let ([u (fv)])
`(,r ,d (lambda (,u) ,(ctx u))))])))))]))])
(cps1 exp id))))%
```


* pmatch

```scheme
(define-syntax pmatch
  (syntax-rules (else guard)
    ((_ (rator rand ...) cs ...)
     (let ((v (rator rand ...)))
       (pmatch v cs ...)))
    ((_ v) (error 'pmatch "failed: ~s" v))
    ((_ v (else e0 e ...)) (begin e0 e ...))
    ((_ v (pat (guard g ...) e0 e ...) cs ...)
     (let ((fk (lambda () (pmatch v cs ...))))
       (ppat v pat (if (and g ...) (begin e0 e ...) (fk)) (fk))))
    ((_ v (pat e0 e ...) cs ...)
     (let ((fk (lambda () (pmatch v cs ...))))
       (ppat v pat (begin e0 e ...) (fk))))))

(define-syntax ppat
  (syntax-rules (_ quote unquote)
    ((_ v _ kt kf) kt)
    ((_ v () kt kf) (if (null? v) kt kf))
    ((_ v (quote lit) kt kf) (if (equal? v (quote lit)) kt kf))
    ((_ v (unquote var) kt kf) (let ((var v)) kt))
    ((_ v (x . y) kt kf)
     (if (pair? v)
       (let ((vx (car v)) (vy (cdr v)))
	 (ppat vx x (ppat vy y kt kf) kf))
       kf))
    ((_ v lit kt kf) (if (equal? v (quote lit)) kt kf))))

```


## test

```scheme
;;; tests

;; var
(cps 'x)
(cps '(lambda (x) x))
(cps '(lambda (x) (x 1)))


;; no lambda (will generate identity functions to return to the toplevel)
(cps '(if (f x) a b))
(cps '(if x (f a) b))


;; if stand-alone (tail)
(cps '(lambda (x) (if (f x) a b)))


;; if inside if-test (non-tail)
(cps '(lambda (x) (if (if x (f a) b) c d)))


;; both branches are trivial, should do some more optimizations
(cps '(lambda (x) (if (if x (zero? a) b) c d)))


;; if inside if-branch (tail)
(cps '(lambda (x) (if t (if x (f a) b) c)))


;; if inside if-branch, but again inside another if-test (non-tail)
(cps '(lambda (x) (if (if t (if x (f a) b) c) e w)))


;; if as operand (non-tail)
(cps '(lambda (x) (h (if x (f a) b))))


;; if as operator (non-tail)
(cps '(lambda (x) ((if x (f g) h) c)))


;; why we need more than two names
(cps '(((f a) (g b)) ((f c) (g d))))



;; factorial
(define fact-cps
(cps
'(lambda (n)
((lambda (fact)
((fact fact) n))
(lambda (fact)
(lambda (n)
(if (zero? n)
1
(* n ((fact fact) (sub1 n))))))))))

;; print out CPSed function
(pretty-print fact-cps)
;; =>
;; '(lambda (n k)
;; ((lambda (fact k) (fact fact (lambda (v0) (v0 n k))))
;; (lambda (fact k)
;; (k
;; (lambda (n k)
;; (if (zero? n)
;; (k 1)
;; (fact
;; fact
;; (lambda (v1) (v1 (sub1 n) (lambda (v2) (k (* n v2))))))))))
;; k))


((eval fact-cps) 5 (lambda (v) v))
;; => 120
```


##  大佬的代码解析

![](http://112firshme11224.test.upcdn.net/blog/tmp/cps.jpg)

![](http://112firshme11224.test.upcdn.net/blog/tmp/cps-1.png)



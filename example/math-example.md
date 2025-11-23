---
layout: math
title: Math Rendering Example
---

# Math Rendering with gojekyll

This page demonstrates mathematical expression rendering using MathJax or KaTeX.

## Inline Math

The famous equation $$E=mc^2$$ shows the relationship between energy and mass.

Here's another example: $$x_0 = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}$$

## Display Math

The Fundamental Theorem of Calculus:

$$
\int_a^b f'(x) dx = f(b) - f(a)
$$

Maxwell's Equations:

$$
\begin{aligned}
\nabla \cdot \mathbf{E} &= \frac{\rho}{\epsilon_0} \\
\nabla \cdot \mathbf{B} &= 0 \\
\nabla \times \mathbf{E} &= -\frac{\partial \mathbf{B}}{\partial t} \\
\nabla \times \mathbf{B} &= \mu_0\left(\mathbf{J} + \epsilon_0 \frac{\partial \mathbf{E}}{\partial t}\right)
\end{aligned}
$$

## Matrices

$$
\begin{bmatrix}
a & b \\
c & d
\end{bmatrix}
\begin{bmatrix}
x \\
y
\end{bmatrix}
=
\begin{bmatrix}
ax + by \\
cx + dy
\end{bmatrix}
$$

## Complex Expression

The probability density function of the normal distribution:

$$
f(x | \mu, \sigma^2) = \frac{1}{\sqrt{2\pi\sigma^2}} e^{-\frac{(x-\mu)^2}{2\sigma^2}}
$$

## Note

Math expressions use the `$$...$$` delimiter for both inline and display math.
- **Inline**: `$$E=mc^2$$` renders as $$E=mc^2$$
- **Display**: Use `$$...$$` on its own lines for centered equations

The delimiters are preserved in the HTML and rendered by MathJax or KaTeX on the client side.

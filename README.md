# ⚙️ SDV_CLI-preprod

Interface en ligne de commande (CLI) permettant la gestion de machines virtuelles dans un environnement VMware vCenter, via son API REST.

---

## 📚 Sommaire

- [📘 Présentation](#-présentation)
- [🧱 Architecture générale](#-architecture-générale)
- [⚙️ Installation](#-installation)
- [🚀 Commandes disponibles](#-commandes-disponibles)
- [🧠 Détails techniques](#-détails-techniques)
- [🛠️ Exemples d’utilisation](#-exemples-dutilisation)
- [📌 Pistes d’amélioration](#-pistes-damélioration)
- [👥 Auteurs](#-auteurs)

---

## 📘 Présentation

`SDV_CLI` est un projet écrit en Go qui propose une interface CLI pour interagir avec une infrastructure VMware (vCenter). Il permet aux utilisateurs de :

- S’authentifier automatiquement via l’API REST
- Gérer leurs VMs (liste, création, démarrage, arrêt, suppression)

---

## 🧱 Architecture générale

```
SDV_CLI-preprod/
├── main.go             # Point d'entrée principal
├── Auth/               # Gestion de l'authentification API
├── Database/           # Stockage local des VM créées
├── cmd/                # Toutes les commandes CLI
```

---

## ⚙️ Installation

```bash
git clone <repo>
cd SDV_CLI-preprod
go mod tidy
go build -o SDVCLI
```

---

## 🚀 Commandes disponibles

| Commande                | Description                                    |
|-------------------------|------------------------------------------------|
| `vm list`               | Liste toutes les VMs                          |
| `vm show <id>`          | Affiche les détails d’une VM                  |
| `vm create`             | Crée une VM (avec saisie interactive)         |
| `vm start <id>`         | Démarre une VM                                |
| `vm stop <id>`          | Arrête une VM                                 |
| `vm delete <id>`        | Supprime une VM                               |

---

## 🧠 Détails techniques

### 🔐 Authentification (`Auth/main.go`)

Ce module gère la session avec l’API REST de vCenter :

- Envoie une requête `POST` vers `/rest/com/vmware/cis/session` avec `BasicAuth`
- Récupère un token (`sessionID`) en cas de succès
- Retourne ce token pour les futurs appels
- Le token est utilisé dans l’en-tête `vmware-api-session-id` des requêtes

### 🧩 Commandes CLI (`cmd/`)

#### `cmd/root.go`

- Point d’entrée de l’application CLI
- Définit la commande `root`, appelée implicitement via `./SDVCLI`
- Initialise les sous-commandes

#### `cmd/vm.go`

- Contient la commande `vm` (groupe logique)
- Attache les sous-commandes : `create`, `start`, `stop`, `show`, `list`, `delete`


#### `cmd/vmCreate.go`

- Permet de créer une VM
- Invite l’utilisateur à saisir :
  - nom de la VM
  - nombre de vCPU
  - quantité de RAM
- Construit un payload JSON
- Envoie une requête `POST` à l’API vSphere `/rest/vcenter/vm`

##### 🎛️ Arguments disponibles

| Argument             | Raccourci | Type   | Description                                      |
|----------------------|-----------|--------|--------------------------------------------------|
| `--name`             | `-n`      | string | Nom de la VM à créer                             |
| `--guest-os`         | `-g`      | string | Type du système invité (ex : `DEBIAN_10_64`)     |
#### `cmd/vmList.go`

- Récupère la liste des VMs via `GET /rest/vcenter/vm`
- Affiche ID, nom, état d’alimentation et hôte

#### `cmd/vmShow.go`

- Utilise `GET /rest/vcenter/vm/<ID>`
- Affiche des informations détaillées sur la VM (CPU, RAM, OS, état)

#### `cmd/vm_start.go` et `cmd/vm_stop.go`

- Envoient respectivement une requête `POST` vers :
  - `/rest/vcenter/vm/<ID>/power/start`
  - `/rest/vcenter/vm/<ID>/power/stop`

#### `cmd/delete.go`

- Envoie une requête `DELETE /rest/vcenter/vm/<ID>`
- Supprime la VM ciblée par ID

---

## 🛠️ Exemples d’utilisation

### Lister les VMs

```bash
./SDVCLI vm list
```

### Démarrer une VM

```bash
./SDVCLI vm start vm-42
```

### Créer une nouvelle VM

```bash
./SDVCLI vm create --name "test" --guest-os "DEBIAN_10_64"
# Invite l'utilisateur à saisir les infos de la VM
```

---

## 📌 Pistes d’amélioration

-Créer des volumes, et les attacher à l'instance
-pouvoir gérer les interfaces réseaux
-Gestion des droits dans une organisation
-Mis en place de la CI/CD

---

## 👥 Auteurs

Projet réalisé dans le cadre du projet d'étude M2 DevOps.

- **Étudiant.e.s** : *Yasmine ABBAS*, *Thibault JUBERT*, *Adem BARAN*
- **Encadrant** : *Morgan KLEIN*

---

## 📄 Licence

Projet pédagogique — tous droits réservés dans le cadre de l’établissement de formation.

const express = require('express');
const connection = require('./database');
const app = express();
const port = 3000;

app.use(express.json());

// Endpoint pour récupérer tous les employés
app.get('/employes', (req, res) => {
    connection.query('SELECT * FROM employes', (error, results) => {
        if (error) {
            return res.status(500).send('Erreur lors de la récupération des employés');
        }
        res.json(results);
    });
});

// Endpoint pour ajouter un employé
app.post('/employes', (req, res) => {
    const { nom, prenom, email, telephone, date_embauche, id_departement, id_poste } = req.body;
    const sql = 'INSERT INTO employes (nom, prenom, email, telephone, date_embauche, id_departement, id_poste) VALUES (?, ?, ?, ?, ?, ?, ?)';
    connection.query(sql, [nom, prenom, email, telephone, date_embauche, id_departement, id_poste], (error, results) => {
        if (error) {
            return res.status(500).send('Erreur lors de l\'ajout de l\'employé');
        }
        res.status(201).send('Employé ajouté avec succès');
    });
});

app.listen(port, () => {
    console.log(`Serveur démarré sur le port ${port}`);
});
